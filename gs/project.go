package gs

import (
	"fmt"
	"github.com/leighmacdonald/verimapcom/pb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HotSpot struct {
	ID    int64
	Lat   float64
	Lon   float64
	Delta float64
}

type Position struct {
	At        time.Time
	Lat       float64
	Lon       float64
	Elevation float64
}

type Project struct {
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	ProjectID    int32  `yaml:"project_id"`
	hotspotsIn   chan pb.HotSpotEvent
	positionsIn  chan pb.PositionEvent
	hotspotFile  *os.File
	positionFile *os.File
}

func (p *Project) AddPosition(pos Position) {
	row := fmt.Sprintf("%d,%.5f,%.5f,%.2f\n", pos.At.Unix(), pos.Lat, pos.Lon, pos.Elevation)
	if _, err := p.positionFile.WriteString(row); err != nil {
		log.Errorf("Could not write position row: %s", row)
	}
}

func (p *Project) AddHotspot(hs HotSpot) {
	row := fmt.Sprintf("%d,%.5f,%.5f,%.2f\n", hs.ID, hs.Lat, hs.Lon, hs.Delta)
	if _, err := p.hotspotFile.WriteString(row); err != nil {
		log.Errorf("Could not write hotspot row: %s", row)
	}
}

func OpenProject(basePath string, projectID int32) (*Project, error) {
	var maxId int64
	found := false
	err := filepath.Walk(basePath, func(p string, info os.FileInfo, err error) error {
		if IsDir(path.Join(basePath, p)) {
			if strings.HasPrefix(p, "project_") {
				pcs := strings.Split(p, "_")
				if len(pcs) != 2 {
					return nil
				}
				pid, err := strconv.ParseInt(pcs[1], 10, 32)
				if err != nil {
					return nil
				}
				if projectID > 0 {
					found = true
				}
				if pid > maxId {
					maxId = pid
				}
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse project directory")
	}
	if !found {
		projectID = int32(maxId) + 1
		pp := fmt.Sprintf("project_%d/ir_export", projectID)
		if err := os.MkdirAll(pp, 0766); err != nil {
			return nil, errors.Wrapf(err, "Could not create project dir: %s", pp)
		}
	}
	p, err := NewProject(projectID)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create project: %s", err)
	}
	return p, nil
}

func NewProject(projectID int32) (*Project, error) {
	basePath := fmt.Sprintf("project_%d", projectID)
	hsPath := path.Join(basePath, "hotspots.csv")
	hsFile, err := os.Create(hsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open new hotspots file: %s", hsPath)
	}
	posPath := path.Join(basePath, "positions.csv")
	posFile, err := os.Create(posPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open new positions file: %s", posPath))
	}
	return &Project{
		Path:         basePath,
		ProjectID:    projectID,
		hotspotsIn:   make(chan pb.HotSpotEvent),
		positionsIn:  make(chan pb.PositionEvent),
		positionFile: posFile,
		hotspotFile:  hsFile,
	}, nil
}
