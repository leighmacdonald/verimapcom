package store

import (
	"github.com/jackc/pgx/pgtype"
	"sync"
	"time"
)

type Person struct {
	PersonID     int32
	AgencyID     int32
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Deleted      bool
	CreatedOn    time.Time
	Agency       Agency
}

type ContextualPerson struct {
	Person
	MissionID int32
	FlightID  int32
	Mu        *sync.RWMutex
}

type BoundingBox struct {
	LatUL  float64
	LongUL float64
	LatLR  float64
	LongLR float64
}

type MissionState int

const (
	StateDeleted    MissionState = 0
	StateCreated    MissionState = 1
	StateLive       MissionState = 2
	StateProcessing MissionState = 3
	StatePublished  MissionState = 4
)

type Evt int

const (
	EvtConnect    Evt = 1
	EvtPing       Evt = 2
	EvtPong       Evt = 3
	EvtMessage    Evt = 10
	EvtSetMission Evt = 20
	EvtError      Evt = 10000
)

type Mission struct {
	MissionID          int32
	PersonID           int32
	AgencyID           int32
	MissionName        string
	MissionState       int32
	ScheduledStartDate pgtype.Timestamp
	ScheduledEndDate   pgtype.Timestamp
	CreatedOn          time.Time
	UpdatedOn          time.Time
	BoundingBox        BoundingBox
	AgencyName         string
	PersonName         string
}

func NewMission(p Person) Mission {
	var m Mission
	m.MissionState = 1
	m.MissionName = "Unnamed Mission"
	m.AgencyID = p.AgencyID
	m.PersonID = p.PersonID
	m.CreatedOn = time.Now()
	m.UpdatedOn = m.CreatedOn
	return m
}

type Message struct {
	MessageID int
	Author    string
	AuthorID  int
	Subject   string
	Body      string
	MissionID int32
	CreatedOn time.Time
}

type Position struct {
	Lat float64
	Lon float64
}

type PositionZ struct {
	Position
	Elevation int32
	ID        int32
	CreatedOn time.Time
	FlightID  int32
}

type HotSpot struct {
	Position
	ID       int32
	Delta    int32
	FlightID int32
}

type FlightState int

const (
	FlightCreated  FlightState = 1
	FlightAirborne FlightState = 2
	FlightLanded   FlightState = 3
	FlightClosed   FlightState = 10
)

type Flight struct {
	FlightID      int32
	MissionID     int32
	FlightState   FlightState
	EnginesOnTime pgtype.Timestamp
	TakeOffTime   pgtype.Timestamp
	LandingTime   pgtype.Timestamp
	Summary       string
	CreatedOn     time.Time

	//hotspotsIn  chan pb.HotSpotEvent
	//positionsIn chan pb.PositionEvent
}

type File struct {
	FileID    int32
	PersonID  int32
	FileName  string
	FileType  string
	FileSize  int64
	Hash256   []byte
	Data      []byte
	CreatedOn time.Time
	UpdatedOn time.Time
}
