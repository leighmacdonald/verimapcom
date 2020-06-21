package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

func MustConnectDB(ctx context.Context) *pgxpool.Pool {
	conn, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}
