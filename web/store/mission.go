package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

type Mission struct {
	MissionID    int
	PersonID     int
	AgencyID     int
	MissionName  string
	MissionState MissionState
	CreatedOn    time.Time
	UpdatedOn    time.Time

	AgencyName string
	PersonName string
}

func (m *Mission) NewFlight() Flight {
	return Flight{
		MissionID:   m.MissionID,
		FlightState: FlightCreated,
		CreatedOn:   time.Now(),
	}
}

func LoadMission(ctx context.Context, db *pgxpool.Pool, missionID int, m *Mission) error {
	const q = `
		SELECT 
		    m.mission_id, m.person_id, m.agency_id, m.mission_name, m.mission_state, m.created_on, m.updated_on,
		       a.agency_name, CONCAT(p.first_name, ' ' , p.last_name) as name
		FROM 
		    mission m
		LEFT JOIN agency a on a.agency_id = m.agency_id
		left join person p on m.person_id = p.person_id
		WHERE 
		    m.mission_id = $1`
	if err := db.QueryRow(ctx, q, missionID).
		Scan(&m.MissionID, &m.PersonID, &m.AgencyID, &m.MissionName,
			&m.MissionState, &m.CreatedOn, &m.UpdatedOn,
			&m.AgencyName, &m.PersonName); err != nil {
		return err
	}
	return nil
}

func InsertMission(ctx context.Context, db *pgxpool.Pool, m *Mission) error {
	m.CreatedOn = time.Now()
	m.UpdatedOn = time.Now()
	const q = `
		INSERT INTO mission (
		   person_id, agency_id, mission_name, mission_state, created_on, updated_on
		) VALUES (
		    $1, $2, $3, $4, $5, $6
		) RETURNING mission_id`

	err := db.QueryRow(ctx, q, m.PersonID, m.AgencyID, m.MissionName,
		m.MissionState, m.CreatedOn, m.UpdatedOn).Scan(&m.MissionID)
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
		    person_id = $2, agency_id = $3, mission_name = $4, mission_state = $5, updated_on = $6
		WHERE
			mission_id = $1`

	if _, err := db.Exec(ctx, q, m.MissionID, m.PersonID, m.AgencyID,
		m.MissionName, m.MissionState, m.UpdatedOn); err != nil {
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

func GetMissions(ctx context.Context, db *pgxpool.Pool, agencyID int) ([]Mission, error) {
	const q = `
		SELECT 
		    m.mission_id, m.person_id, m.agency_id, m.mission_name, m.mission_state, m.created_on, m.updated_on,
		    a.agency_name, CONCAT(p.first_name, ' ' , p.last_name) as name
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
			&m.MissionState, &m.CreatedOn, &m.UpdatedOn, &m.AgencyName, &m.PersonName); err != nil {
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
