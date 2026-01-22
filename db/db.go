// Package db
package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func getURL() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environment variables: %v\n", err)
	}

	url := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	return url
}

func NewDBPool() (*pgxpool.Pool, error) {
	url := getURL()
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create connection pool: %v\n", err)
	}
	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatal("Couldnt connect to database: ", err)
	}
	return dbpool, nil
}
