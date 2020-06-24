package store

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

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

type BoundingBox struct {
	LatUL  float64
	LongUL float64
	LatLR  float64
	LongLR float64
}

type Mission struct {
	MissionID          int
	PersonID           int
	AgencyID           int
	MissionName        string
	MissionState       int
	ScheduledStartDate pgtype.Timestamp
	ScheduledEndDate   pgtype.Timestamp
	CreatedOn          time.Time
	UpdatedOn          time.Time
	BoundingBox        BoundingBox
	AgencyName         string
	PersonName         string
}

func (m *Mission) NewFlight() Flight {
	return Flight{
		MissionID:   m.MissionID,
		FlightState: FlightCreated,
		CreatedOn:   time.Now(),
	}
}

func InsertMission(ctx context.Context, db *pgxpool.Pool, m *Mission) error {
	m.CreatedOn = time.Now()
	m.UpdatedOn = time.Now()
	const q = `
		INSERT INTO mission (
		   person_id, agency_id, mission_name, mission_state, created_on, updated_on,
			bbox_ul, bbox_lr
		) VALUES (
		    $1, $2, $3, $4, $5, $6, 
		      CONCAT('SRID=4326;POINT(', $8::float8, ' ', $7::float8, ')'),
		    CONCAT('SRID=4326;POINT(', $10::float8,' ', $9::float8, ')')
		) RETURNING mission_id`

	err := db.QueryRow(ctx, q, m.PersonID, m.AgencyID, m.MissionName,
		m.MissionState, m.CreatedOn, m.UpdatedOn, m.BoundingBox.LatUL, m.BoundingBox.LongUL,
		m.BoundingBox.LatLR, m.BoundingBox.LongLR).Scan(&m.MissionID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMission(ctx context.Context, db *pgxpool.Pool, m *Mission) error {
	m.UpdatedOn = time.Now()
	const q = `
		UPDATE 
		    mission
		SET 
		    person_id = $2, agency_id = $3, mission_name = $4, mission_state = $5, updated_on = $6,
			bbox_ul = CONCAT('SRID=4326;POINT(', $8, ' ', $7, ')'),
		    bbox_lr = CONCAT('SRID=4326;POINT(', $10,' ', $9, ')'), 
		    scheduled_start_date = $11, scheduled_end_date = $12
		WHERE
			mission_id = $1`

	if _, err := db.Exec(ctx, q, m.MissionID, m.PersonID, m.AgencyID,
		m.MissionName, m.MissionState, m.UpdatedOn, m.BoundingBox.LatUL, m.BoundingBox.LongUL,
		m.BoundingBox.LatLR, m.BoundingBox.LongLR, m.ScheduledStartDate,
		m.ScheduledEndDate); err != nil {
		return err
	}
	return nil
}

func SaveMission(ctx context.Context, db *pgxpool.Pool, m *Mission) error {
	if m.MissionID > 0 {
		return UpdateMission(ctx, db, m)
	}
	return InsertMission(ctx, db, m)
}

func GetMission(ctx context.Context, db *pgxpool.Pool, missionID int) (Mission, error) {
	const q = `
		SELECT 
		    m.mission_id, m.person_id, m.agency_id, m.mission_name, m.mission_state, m.created_on, m.updated_on,
		    a.agency_name, CONCAT(p.first_name, ' ' , p.last_name) as name,
		    ST_Y(bbox_ul), ST_X(bbox_ul), ST_Y(bbox_lr),ST_X(bbox_lr), m.scheduled_start_date, m.scheduled_end_date
		FROM 
		    mission m
		LEFT JOIN agency a on a.agency_id = m.agency_id
		LEFT JOIN person p on m.person_id = p.person_id
		WHERE m.mission_id = $1`
	var m Mission
	if err := db.QueryRow(ctx, q, missionID).Scan(&m.MissionID, &m.PersonID, &m.AgencyID, &m.MissionName,
		&m.MissionState, &m.CreatedOn, &m.UpdatedOn, &m.AgencyName, &m.PersonName, &m.BoundingBox.LatUL, &m.BoundingBox.LongUL,
		&m.BoundingBox.LatLR, &m.BoundingBox.LongLR, &m.ScheduledStartDate,
		&m.ScheduledEndDate); err != nil {
		return m, err
	}
	return m, nil
}

func GetMissions(ctx context.Context, db *pgxpool.Pool, agencyID int) ([]Mission, error) {
	const q = `
		SELECT 
		    m.mission_id, m.person_id, m.agency_id, m.mission_name, m.mission_state, m.created_on, m.updated_on,
		    a.agency_name, CONCAT(p.first_name, ' ' , p.last_name) as name,
		    ST_Y(bbox_ul), ST_X(bbox_ul), ST_Y(bbox_lr),ST_X(bbox_lr), m.scheduled_start_date, m.scheduled_end_date
		FROM 
		    mission m
		LEFT JOIN agency a on a.agency_id = m.agency_id
		LEFT JOIN person p on m.person_id = p.person_id`
	query := q
	var (
		missions []Mission
		rows     pgx.Rows
		err      error
	)
	var args []interface{}
	if agencyID > 1 {
		query += " WHERE agency_id = $1"
		args = append(args, agencyID)
	}
	rows, err = db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Mission
		if err := rows.Scan(&m.MissionID, &m.PersonID, &m.AgencyID, &m.MissionName,
			&m.MissionState, &m.CreatedOn, &m.UpdatedOn, &m.AgencyName, &m.PersonName,
			&m.BoundingBox.LatUL, &m.BoundingBox.LongUL,
			&m.BoundingBox.LatLR, &m.BoundingBox.LongLR,
			&m.ScheduledStartDate, &m.ScheduledEndDate); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func MissionAttachFile(ctx context.Context, db *pgxpool.Pool, missionID int, fileID int) error {
	const q = `
		INSERT INTO mission_file (
		    mission_id, file_id, created_on
		) VALUES ($1, $2, $3)`
	if _, err := db.Exec(ctx, q, missionID, fileID, time.Now()); err != nil {
		return err
	}
	return nil
}

func MissionDetachFile(ctx context.Context, db *pgxpool.Pool, missionID int, fileID int) error {
	const q = `
		DELETE FROM  mission_file WHERE mission_id = $1 AND file_id = $2`
	if _, err := db.Exec(ctx, q, missionID, fileID); err != nil {
		return err
	}
	return nil
}

type MissionEvent struct {
	MissionEventID int                    `json:"mission_event_id"`
	MissionID      int                    `json:"mission_id"`
	EventType      Evt                    `json:"event_type"`
	Payload        map[string]interface{} `json:"payload"`
	CreatedOn      time.Time              `json:"created_on"`
}

func NewMissionEvent(evt Evt, missionID int) MissionEvent {
	return MissionEvent{
		MissionID: missionID,
		EventType: evt,
		Payload:   make(map[string]interface{}),
		CreatedOn: time.Now(),
	}
}

func MissionEventAdd(ctx context.Context, db *pgxpool.Pool, e *MissionEvent) error {
	const q = `
		INSERT INTO mission_event 
		    (mission_id, event_type, payload, created_on) 
		VALUES 
		    ($1, $2, $3, $4) 
		RETURNING 
		    mission_event_id`
	if err := db.QueryRow(ctx, q, e.MissionID, e.EventType, e.Payload, e.CreatedOn).Scan(&e.MissionEventID); err != nil {
		return errors.Wrapf(err, "failed to insert mission event")
	}
	return nil
}

func MissionEventGetAll(ctx context.Context, db *pgxpool.Pool, missionID int) ([]MissionEvent, error) {
	const q = `
		SELECT mission_event_id, mission_id, event_type, payload, created_on 
		FROM mission_event
		WHERE mission_id = $1`
	rows, err := db.Query(ctx, q, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []MissionEvent
	for rows.Next() {
		var e MissionEvent
		if err := rows.Scan(&e.MissionEventID, &e.MissionID, &e.EventType, &e.Payload, &e.CreatedOn); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
