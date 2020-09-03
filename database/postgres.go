package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DBConn - connects to postgresql database
func DBConn() (*pgxpool.Pool, error) {
	// connecting to databbase
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	// error handling
	if err != nil {
		return nil, err
	}

	log.Println("connecting to dattabase...")

	return conn, err
}
