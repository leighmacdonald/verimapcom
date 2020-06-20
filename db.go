package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

func mustConnectDB(ctx context.Context) *pgxpool.Pool {
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}

type Agency struct {
	AgencyID   int
	AgencyName string
	CreatedOn  time.Time
}

type Person struct {
	PersonID     int
	AgencyID     int
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedOn    time.Time
	Agency       Agency
}

func getPeople() ([]Person, error) {
	var people []Person
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id`
	rows, err := dbpool.Query(ctx, q)
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

func loadPersonByEmail(email string, person *Person) error {
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id
		WHERE 
		    email = $1`
	if err := dbpool.QueryRow(ctx, q, email).
		Scan(&person.PersonID, &person.AgencyID, &person.Email, &person.PasswordHash,
			&person.FirstName, &person.LastName, &person.CreatedOn,
			&person.Agency.AgencyID, &person.Agency.AgencyName, &person.Agency.CreatedOn); err != nil {
		return err
	}
	return nil
}

func loadPersonByID(userID int, person *Person) error {
	const q = `
		SELECT 
			p.person_id, p.agency_id, p.email, p.password_hash, p.first_name, p.last_name, p.created_on, 
		       a.agency_id, a.agency_name, a.created_on
		FROM 
		    person p
		LEFT JOIN agency a on a.agency_id = p.agency_id
		WHERE 
		    person_id = $1`
	if err := dbpool.QueryRow(ctx, q, userID).
		Scan(&person.PersonID, &person.AgencyID, &person.Email, &person.PasswordHash,
			&person.FirstName, &person.LastName, &person.CreatedOn,
			&person.Agency.AgencyID, &person.Agency.AgencyName, &person.Agency.CreatedOn); err != nil {
		return err
	}
	return nil
}

func insertPerson(person *Person) error {
	const q = `
		INSERT INTO person (
		    agency_id, email, password_hash, first_name, last_name, created_on)  
		VALUES (
		    $1, $2, $3, $4, $5, $6
		) RETURNING person_id`

	err := dbpool.QueryRow(ctx, q, person.AgencyID, person.Email, person.PasswordHash,
		person.FirstName, person.LastName, person.CreatedOn).Scan(&person.PersonID)
	if err != nil {
		return err
	}
	return nil
}

func updatePerson(person *Person) error {
	const q = `
		UPDATE 
		    person
		SET 
		    agency_id=$2, email=$3, password_hash=$4, first_name = $5, last_name = $6, created_on = $7  
		WHERE
			person_id = $1`

	if _, err := dbpool.Exec(ctx, q, person.PersonID, person.AgencyID, person.Email, person.PasswordHash,
		person.FirstName, person.LastName, person.CreatedOn); err != nil {
		return err
	}
	return nil
}

func savePerson(person *Person) error {
	if person.PersonID > 0 {
		return updatePerson(person)
	}
	return insertPerson(person)
}

func getAgencies() ([]Agency, error) {
	const q = `
		SELECT 
			agency_id, agency_name, created_on 
		FROM 
		    agency`
	var agencies []Agency
	rows, err := dbpool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a Agency
		if err := rows.Scan(&a.AgencyID, &a.AgencyName, &a.CreatedOn); err != nil {
			return nil, err
		}
		agencies = append(agencies, a)
	}
	return agencies, nil
}

func loadAgency(agencyID int, agency *Agency) error {
	const q = `
		SELECT 
			agency_id, agency_name, created_on 
		FROM 
		    agency 
		WHERE 
		    agency_id = $1`
	if err := dbpool.QueryRow(ctx, q, agencyID).
		Scan(&agency.AgencyID, &agency.AgencyName, &agency.CreatedOn); err != nil {
		return err
	}
	return nil
}

func insertAgency(agency *Agency) error {
	const q = `
		INSERT INTO agency (
		   agency_name, created_on 
		) VALUES (
		    $1, $2
		) RETURNING agency_id`

	err := dbpool.QueryRow(ctx, q, agency.AgencyName, agency.CreatedOn).Scan(&agency.AgencyID)
	if err != nil {
		return err
	}
	return nil
}

func updateAgency(person *Agency) error {
	const q = `
		UPDATE 
		    agency
		SET 
		    agency_name = $2, created_on = $3  
		WHERE
			agency_id = $1`

	if _, err := dbpool.Exec(ctx, q, person.AgencyID, person.CreatedOn, person.CreatedOn); err != nil {
		return err
	}
	return nil
}

func saveAgency(agency *Agency) error {
	if agency.AgencyID > 0 {
		return updateAgency(agency)
	}
	return insertAgency(agency)
}
