package cmd

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/verimapcom/core"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
		ctx, cancel := context.WithCancel(context.Background())
		irSrc, err := cmd.Flags().GetString("ir")
		if err != nil {
			log.Fatalf("--ir must be set")
		}
		if !core.Exists(irSrc) {
			log.Fatalf("--ir must exist (%s)", irSrc)
		}
		dst, err := cmd.Flags().GetString("dest")
		if err != nil {
			log.Fatalf("--dest must be set")
		}
		if core.Exists(dst) {
			if err := os.RemoveAll(dst); err != nil {
				log.Fatalf("Error trying to remove scratch dir: %v", err)
			}
		}
		irOutPath := path.Join(dst, "ir_export")
		if err := os.MkdirAll(irOutPath, 0766); err != nil {
			log.Fatalf("Error trying to create scratch dir")
		}

		position, err := cmd.Flags().GetString("position")
		if err != nil {
			log.Fatalf("--position must be set")
		}
		if !core.Exists(position) {
			log.Fatalf("--position must exist (%s)", position)
		}

		hotspots, err := cmd.Flags().GetString("hotspots")
		if err != nil {
			log.Fatalf("--hotspots must be set")
		}
		if !core.Exists(hotspots) {
			log.Fatalf("--hotspots must exist (%s)", hotspots)
		}
		var fileNames []string
		rxPNGExport := regexp.MustCompile(`^.+?_\d+_\d+_level_1.(png|wld)`)
		if err := filepath.Walk(irSrc, func(path string, info os.FileInfo, err error) error {
			name := filepath.Base(path)
			if rxPNGExport.MatchString(name) && !info.IsDir() {
				fileNames = append(fileNames, name)
			}
			return nil
		}); err != nil {
			log.Fatalf("Failed to walk directory: %v", err)
		}
		go func() {
			i := 0
			t := time.NewTicker(time.Second * 5)
			dstDir := path.Join(dst, "ir_export")
			if err := os.MkdirAll(dstDir, 0766); err != nil {
				log.Fatalf("Failed to create ir_export dir: %v", err)
			}
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if i+1 > len(fileNames) {
						log.Info("Export IR queue complete")
						return
					}
					input, err := ioutil.ReadFile(path.Join(irSrc, fileNames[i]))
					if err != nil {
						fmt.Println(err)
						return
					}
					err = ioutil.WriteFile(path.Join(dstDir, filepath.Base(fileNames[i])), input, 0644)
					if err != nil {
						log.Fatalf("Error creating: %s", filepath.Base(fileNames[i]))
						return
					}
					log.Infof("Wrote frame data %d: %s", i, path.Join(dstDir, filepath.Base(fileNames[i])))
					i++
				}
			}
		}()
		go func() {
			input, err := ioutil.ReadFile(position)
			if err != nil {
				log.Fatalf(fmt.Sprintf("failed to read positions file: %v", err))
				return
			}
			outFile := path.Join(dst, "stage_position_file.csv")
			fp, err := os.Create(outFile)
			if err != nil {
				log.Fatalf("Error creating %s: %v", outFile, err)
			}
			defer func() {
				if err := fp.Close(); err != nil {
					log.Errorf("Failed to cleanly close position file: %v", err)
				}
			}()
			rows := strings.Split(string(input), "\r\n")
			i := 0
			t := time.NewTicker(time.Millisecond * 10)
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if i+1 > len(rows) {
						log.Info("position queue complete")
						return
					}
					// do stuff
					if _, err := fp.WriteString(rows[i] + "\r\n"); err != nil {
						log.Fatalf("Error writing position line: %v", err)
					}
					if i%1000 == 0 {
						log.Infof("Wrote Pos %d: %s", i, outFile)
					}
					i++
				}
			}
		}()
		go func() {
			input, err := ioutil.ReadFile(hotspots)
			if err != nil {
				log.Fatalf(fmt.Sprintf("failed to read hotspots file: %v", err))
				return
			}
			hotspotFile := path.Join(dst, "hotcluster.txt")
			fp, err := os.Create(hotspotFile)
			if err != nil {
				log.Fatalf("Failed to create hotcluster output file: %v", err)
			}
			defer func() {
				if err := fp.Close(); err != nil {
					log.Errorf("Failed to cleanly close hotcluster")
				}
			}()
			rows := strings.Split(string(input), "\r\n")
			i := 0
			t := time.NewTicker(time.Second * 1)
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if i+1 > len(rows) {
						log.Info("Hotspots queue complete")
						return
					}
					// do stuff
					if _, err := fp.WriteString(rows[i] + "\r\n"); err != nil {
						log.Errorf("Failed to write hotspot row: %v", err)
					}
					if i%10 == 0 {
						log.Infof("Wrote Hotspots: %d: %s", i, hotspotFile)
					}
					i++
				}
			}
		}()
		core.WaitForSignal(ctx, func(ctx context.Context) error {
			cancel()
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)

	demoCmd.Flags().StringP("dest", "D", "", "Where to write the files to monitor")
	demoCmd.Flags().StringP("ir", "s", "", "Path to the export output directory")
	demoCmd.Flags().StringP("position", "P", "stage_position_file.csv", "Path to stage_position")
	demoCmd.Flags().StringP("hotspots", "t", "hotcluster.txt", "Path to hotcluster txt")
}
