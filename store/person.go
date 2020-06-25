package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func PersonDelete(ctx context.Context, db *pgxpool.Pool, personID int32) error {
	const q = `UPDATE person SET deleted = true WHERE person_id = $1`
	if _, err := db.Exec(ctx, q, personID); err != nil {
		return err
	}
	return nil
}

func GetPeople(ctx context.Context, db *pgxpool.Pool) ([]Person, error) {
	var people []Person
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id
		WHERE p.deleted = false
		ORDER BY p.person_id`
	rows, err := db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.PersonID, &p.AgencyID, &p.Email, &p.PasswordHash,
			&p.FirstName, &p.LastName, &p.CreatedOn,
			&p.Agency.AgencyID, &p.Agency.AgencyName, &p.Agency.CreatedOn); err != nil {
			return nil, err
		}
		people = append(people, p)
	}
	return people, nil
}

func LoadPersonByToken(ctx context.Context, db *pgxpool.Pool, rpcToken string, person *Person) error {
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id
		WHERE 
		      p.deleted = false AND rpc_token = $1
		ORDER BY p.person_id`
	if err := db.QueryRow(ctx, q, rpcToken).
		Scan(&person.PersonID, &person.AgencyID, &person.Email, &person.PasswordHash,
			&person.FirstName, &person.LastName, &person.CreatedOn,
			&person.Agency.AgencyID, &person.Agency.AgencyName, &person.Agency.CreatedOn); err != nil {
		return err
	}
	return nil
}

func LoadPersonByEmail(ctx context.Context, db *pgxpool.Pool, email string, person *Person) error {
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id
		WHERE 
		      p.deleted = false AND email = $1
		ORDER BY p.person_id`
	if err := db.QueryRow(ctx, q, email).
		Scan(&person.PersonID, &person.AgencyID, &person.Email, &person.PasswordHash,
			&person.FirstName, &person.LastName, &person.CreatedOn,
			&person.Agency.AgencyID, &person.Agency.AgencyName, &person.Agency.CreatedOn); err != nil {
		return err
	}
	return nil
}

func LoadPersonByID(ctx context.Context, db *pgxpool.Pool, userID int32, p *Person) error {
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on, p.deleted
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id
		WHERE 
		    p.person_id = $1`
	if err := db.QueryRow(ctx, q, userID).
		Scan(&p.PersonID, &p.AgencyID, &p.Email, &p.PasswordHash,
			&p.FirstName, &p.LastName, &p.CreatedOn,
			&p.Agency.AgencyID, &p.Agency.AgencyName, &p.Agency.CreatedOn, &p.Deleted); err != nil {
		return err
	}
	return nil
}

func InsertPerson(ctx context.Context, db *pgxpool.Pool, person *Person) error {
	const q = `
		INSERT INTO person (
		    agency_id, email, password_hash, first_name, last_name, created_on)  
		VALUES (
		    $1, $2, $3, $4, $5, $6
		) RETURNING person_id`

	err := db.QueryRow(ctx, q, person.AgencyID, person.Email, person.PasswordHash,
		person.FirstName, person.LastName, person.CreatedOn).Scan(&person.PersonID)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePerson(ctx context.Context, db *pgxpool.Pool, person *Person) error {
	const q = `
		UPDATE 
		    person
		SET 
		    agency_id=$2, email=$3, password_hash=$4, first_name = $5, last_name = $6, created_on = $7,
			deleted = $8
		WHERE
			person_id = $1`

	if _, err := db.Exec(ctx, q, person.PersonID, person.AgencyID, person.Email, person.PasswordHash,
		person.FirstName, person.LastName, person.CreatedOn, person.Deleted); err != nil {
		return err
	}
	return nil
}

func SavePerson(ctx context.Context, db *pgxpool.Pool, person *Person) error {
	if person.PersonID > 0 {
		return UpdatePerson(ctx, db, person)
	}
	return InsertPerson(ctx, db, person)
}
