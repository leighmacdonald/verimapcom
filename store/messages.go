package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MessageAdd(ctx context.Context, db *pgxpool.Pool, m *Message) error {
	const q = `
		INSERT INTO message 
    		(author_id, author_name, message_subject, message_body, created_on) 
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING 
		    message_id`
	if err := db.QueryRow(ctx, q, m.AuthorID, m.Author,
		m.Subject, m.Body, m.CreatedOn).Scan(&m.MessageID); err != nil {
		return err
	}
	return nil
}
