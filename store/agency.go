package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Agency struct {
	AgencyID   int32
	AgencyName string
	InviteKey  string
	Slots      int
	SlotsUsed  int
	CreatedOn  time.Time
}

func GetAgencies(ctx context.Context, db *pgxpool.Pool) ([]Agency, error) {
	const q = `
		SELECT 
			a.agency_id,
	   		a.agency_name,
	   		a.invite_key,
	   		a.slots,
	   		a.created_on,
	   		COUNT(p.person_id) as slots_used
		FROM agency a
		LEFT JOIN person p on a.agency_id = p.agency_id
		GROUP BY a.agency_id`
	var agencies []Agency
	rows, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a Agency
		if err := rows.Scan(&a.AgencyID, &a.AgencyName, &a.InviteKey, &a.Slots, &a.CreatedOn, &a.SlotsUsed); err != nil {
			return nil, err
		}
		agencies = append(agencies, a)
	}
	return agencies, nil
}

func LoadAgency(ctx context.Context, db *pgxpool.Pool, agencyID int, agency *Agency) error {
	const q = `SELECT 
			a.agency_id,
	   		a.agency_name,
	   		a.invite_key,
	   		a.slots,
	   		a.created_on,
	   		COUNT(p.person_id) as slots_used
		FROM agency a
		LEFT JOIN person p on a.agency_id = p.agency_id
		WHERE a.agency_id = $1
		GROUP BY a.agency_id
		    `
	if err := db.QueryRow(ctx, q, agencyID).
		Scan(&agency.AgencyID, &agency.AgencyName, &agency.InviteKey, &agency.Slots,
			&agency.CreatedOn, &agency.SlotsUsed); err != nil {
		return err
	}
	return nil
}

func LoadAgencyByInviteKey(ctx context.Context, db *pgxpool.Pool, inviteKey string, agency *Agency) error {
	const q = `SELECT 
			a.agency_id,
	   		a.agency_name,
	   		a.invite_key,
	   		a.slots,
	   		a.created_on,
	   		COUNT(p.person_id) as slots_used
		FROM agency a
		LEFT JOIN person p on a.agency_id = p.agency_id
		WHERE a.invite_key = $1
		GROUP BY a.agency_id
		    `
	if err := db.QueryRow(ctx, q, inviteKey).
		Scan(&agency.AgencyID, &agency.AgencyName, &agency.InviteKey, &agency.Slots,
			&agency.CreatedOn, &agency.SlotsUsed); err != nil {
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
