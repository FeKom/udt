package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect reads POSTGRES_URI from the environment and returns a connection pool.
// It pings the database to confirm the connection is live, logging and exiting
// the process if the connection cannot be established.
func Connect() *pgxpool.Pool {
	uri := os.Getenv("POSTGRES_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		log.Fatalf("postgres: failed to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres: failed to ping database: %v", err)
	}

	return pool
}
