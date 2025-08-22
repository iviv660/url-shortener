package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
)

func ConnectPostgres(ctx context.Context, urlPostgres string) (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres")
	return db, nil
}
