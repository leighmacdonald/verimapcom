package store

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type FlightState int

const (
	FlightCreated  FlightState = 1
	FlightAirborne FlightState = 2
	FlightLanded   FlightState = 3
	FlightClosed   FlightState = 10
)

type Flight struct {
	FlightID      int
	MissionID     int
	FlightState   FlightState
	EnginesOnTime pgtype.Timestamp
	TakeOffTime   pgtype.Timestamp
	LandingTime   pgtype.Timestamp
	Summary       string
	CreatedOn     time.Time
}

func FlightSave(ctx context.Context, db *pgxpool.Pool, flight *Flight) error {
	if flight.FlightID > 0 {
		return flightUpdate(ctx, db, flight)
	} else {
		return flightInsert(ctx, db, flight)
	}
}

func Flights(ctx context.Context, db *pgxpool.Pool) ([]Flight, error) {
	var flights []Flight
	const q = `
		SELECT 
			flight_id, mission_id, flight_state, engines_on_time,takeoff_time,
       		landing_time, summary, created_on 
		FROM 
		    flight`
	rows, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f Flight
		if err := rows.Scan(&f.FlightID, &f.MissionID, &f.FlightState, &f.EnginesOnTime,
			&f.TakeOffTime, &f.LandingTime, &f.Summary, &f.CreatedOn); err != nil {
			return nil, err
		}
		flights = append(flights, f)
	}
	return flights, nil
}

func FlightsByMissionID(ctx context.Context, db *pgxpool.Pool, missionID int) ([]Flight, error) {
	var flights []Flight
	const q = `
		SELECT 
			flight_id, mission_id, flight_state, engines_on_time,takeoff_time,
       		landing_time, summary, created_on 
		FROM 
		    flight 
		WHERE 
		    mission_id = $1`
	rows, err := db.Query(ctx, q, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f Flight
		if err := rows.Scan(&f.FlightID, &f.MissionID, &f.FlightState, &f.EnginesOnTime,
			&f.TakeOffTime, &f.LandingTime, &f.Summary, &f.CreatedOn); err != nil {
			return nil, err
		}
		flights = append(flights, f)
	}
	return flights, nil
}

func flightInsert(ctx context.Context, db *pgxpool.Pool, flight *Flight) error {
	const q = `
		INSERT INTO flight (
		    mission_id, created_on, flight_state, summary
		) VALUES ($1, $2) RETURNING flight_id`
	if err := db.QueryRow(ctx, q, flight.MissionID, flight.CreatedOn,
		flight.FlightState, flight.Summary).Scan(&flight.FlightID); err != nil {
		return errors.Wrap(err, "Failed to insert new flight")
	}
	return nil
}

func flightUpdate(ctx context.Context, db *pgxpool.Pool, f *Flight) error {
	const q = `
		UPDATE flight
		SET 
			mission_id = $2, created_on = $3, flight_state = $4, summary =  $5, 
		    engines_on_time = $6, takeoff_time = $7, landing_time = $8
		WHERE	
			flight_id = $1`
	if _, err := db.Exec(ctx, q, f.FlightID, f.MissionID, f.CreatedOn, f.FlightState,
		f.Summary, f.EnginesOnTime, f.TakeOffTime, f.LandingTime); err != nil {
		return errors.Wrap(err, "Failed to update flight")
	}
	return nil
}
