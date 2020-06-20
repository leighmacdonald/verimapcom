package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"time"
)

type MissionState int

const (
	deleted    MissionState = 0
	created    MissionState = 1
	live       MissionState = 2
	processing MissionState = 3
	published  MissionState = 4
)

type Mission struct {
	MissionID    int
	PersonID     int
	AgencyID     int
	MissionName  string
	MissionState MissionState
	CreatedOn    time.Time
	UpdatedOn    time.Time
}

func loadMission(missionID int, m *Mission) error {
	const q = `
		SELECT 
		    mission_id, person_id, agency_id, mission_name, mission_state, created_on, updated_on
		FROM 
		    mission
		WHERE 
		    mission_id = $1`
	if err := dbpool.QueryRow(ctx, q, missionID).
		Scan(&m.MissionID, &m.PersonID, &m.AgencyID, &m.MissionName,
			&m.MissionState, &m.CreatedOn, &m.UpdatedOn); err != nil {
		return err
	}
	return nil
}

func insertMission(m *Mission) error {
	m.CreatedOn = time.Now()
	m.UpdatedOn = time.Now()
	const q = `
		INSERT INTO mission (
		   person_id, agency_id, mission_name, mission_state, created_on, updated_on
		) VALUES (
		    $1, $2, $3, $4, $5, $6
		) RETURNING mission_id`

	err := dbpool.QueryRow(ctx, q, m.PersonID, m.AgencyID, m.MissionName,
		m.MissionState, m.CreatedOn, m.UpdatedOn).Scan(&m.MissionID)
	if err != nil {
		return err
	}
	return nil
}

func updateMission(m *Mission) error {
	m.UpdatedOn = time.Now()
	const q = `
		UPDATE 
		    mission
		SET 
		    person_id = $2, agency_id = $3, mission_name = $4, mission_state = $5, updated_on = $6
		WHERE
			mission_id = $1`

	if _, err := dbpool.Exec(ctx, q, m.MissionID, m.PersonID, m.AgencyID,
		m.MissionName, m.MissionState, m.UpdatedOn); err != nil {
		return err
	}
	return nil
}

func saveMission(m *Mission) error {
	if m.MissionID > 0 {
		return updateMission(m)
	}
	return insertMission(m)
}

func dbGetMissions(agencyID int) ([]Mission, error) {
	const q = `SELECT 
		    mission_id, person_id, agency_id, mission_name, mission_state, created_on, updated_on
		FROM 
		    mission `
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
	rows, err = dbpool.Query(ctx, "", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Mission
		if err := rows.Scan(&m.MissionID, &m.PersonID, &m.AgencyID, &m.MissionName,
			&m.MissionState, &m.CreatedOn, &m.UpdatedOn); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func getMissions(c *gin.Context) {
	m := defaultM(c, missions)
	userMissions, err := dbGetMissions(m["person"].(Person).AgencyID)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		log.Errorf("Failed to fetch user missions: %v", err)
	}
	m["missions"] = userMissions
	render(c, missions, m)
}
