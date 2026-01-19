package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/joho/godotenv"
)

func Greet() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environment variables: %v\n", err)
	}

	url := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
}
