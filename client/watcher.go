package client

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/leighmacdonald/verimapcom/gs"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func eventParser(watcher *fsnotify.Watcher, newFiles chan string) {
	rxPNGExport := regexp.MustCompile(`^.+?_\d+_\d+_level_1.png`)
	posCnt := 0
	hsCnt := 0
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if strings.Contains(event.String(), "stage_position_file") {
					posCnt++
					if posCnt%500 == 0 {
						log.WithFields(log.Fields{
							"event": event.Name,
							"op":    event.Op,
						}).Infof("x%d", 500)
					}
				} else if strings.Contains(event.String(), "hotcluster") {
					hsCnt++
					if posCnt%100 == 0 {
						log.WithFields(log.Fields{
							"event": event.Name,
							"op":    event.Op,
						}).Infof("x%d", 100)
					}
				} else {
					log.Printf("%s", event.Name)
				}
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				name := filepath.Base(event.Name)
				if strings.HasPrefix(name, "hotcluster") {
					newFiles <- event.Name
				} else if name == "stage_position_file.csv" {
					newFiles <- event.Name
				} else {
					if gs.IsDir(event.Name) {
						if err := watcher.Add(event.Name); err != nil {
							log.Errorf("Failed to add new watch dir: %v", err)
							continue
						}
						log.Infof("Added new watch dir: %s", event.Name)
						continue
					}
					fileName := filepath.Base(event.Name)
					if rxPNGExport.MatchString(fileName) {
						log.Infof("Got new png export: %v", fileName)
						newFiles <- event.Name
						continue
					}
					log.Infof("Ignored file created: %s", event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			log.Errorf("error: %s", err)
			if !ok {
				return
			}

		}
	}
}

func monitorDirectory(ctx context.Context, directory string, newFiles chan string) {
	log.Infof("Watching folder: %s", directory)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Errorf("Failed to close watcher cleanly: %s", err)
		}
	}()
	if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Infof("Watching folder: %s", path)
			return watcher.Add(path)
		}
		// TODO dont hardcode
		for _, k := range []string{"hotcluster", "stage_position_file"} {
			if strings.HasPrefix(info.Name(), k) {
				log.Infof("Watching: %s", path)
				newFiles <- path
				return watcher.Add(path)
			}
		}
		return nil
	}); err != nil {
		log.Warnf("Failed to walk subdirs to watch: %v", err)
	}

	go eventParser(watcher, newFiles)
	go func() {
		for {
			select {
			case err := <-watcher.Errors:
				log.Errorf("Watcher err: %v", err)
			}
		}
	}()
	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	log.Infof("Stopped watching dir: %s", directory)
}
