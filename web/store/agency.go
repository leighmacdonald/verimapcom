package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Agency struct {
	AgencyID   int
	AgencyName string
	Slots      int
	CreatedOn  time.Time
}

func GetAgencies(ctx context.Context, db *pgxpool.Pool) ([]Agency, error) {
	const q = `
		SELECT 
			agency_id, agency_name, slots, created_on 
		FROM 
		    agency`
	var agencies []Agency
	rows, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a Agency
		if err := rows.Scan(&a.AgencyID, &a.AgencyName, &a.Slots, &a.CreatedOn); err != nil {
			return nil, err
		}
		agencies = append(agencies, a)
	}
	return agencies, nil
}

func LoadAgency(ctx context.Context, db *pgxpool.Pool, agencyID int, agency *Agency) error {
	const q = `
		SELECT 
			agency_id, agency_name, slots, created_on 
		FROM 
		    agency 
		WHERE 
		    agency_id = $1`
	if err := db.QueryRow(ctx, q, agencyID).
		Scan(&agency.AgencyID, &agency.AgencyName, &agency.Slots, &agency.CreatedOn); err != nil {
		return err
	}
	return nil
}

func InsertAgency(ctx context.Context, db *pgxpool.Pool, agency *Agency) error {
	const q = `
		INSERT INTO agency (
		   agency_name, slots, created_on 
		) VALUES (
		    $1, $2, $3
		) RETURNING agency_id`

	err := db.QueryRow(ctx, q, agency.AgencyName, agency.Slots, agency.CreatedOn).Scan(&agency.AgencyID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAgency(ctx context.Context, db *pgxpool.Pool, a *Agency) error {
	const q = `
		UPDATE 
		    agency
		SET 
		    agency_name = $2, slots = $3, created_on = $4  
		WHERE
			agency_id = $1`

	if _, err := db.Exec(ctx, q, a.AgencyID, a.AgencyName, a.Slots, a.CreatedOn); err != nil {
		return err
	}
	return nil
}

func SaveAgency(ctx context.Context, db *pgxpool.Pool, agency *Agency) error {
	if agency.AgencyID > 0 {
		return UpdateAgency(ctx, db, agency)
	}
	return InsertAgency(ctx, db, agency)
}
