package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgresDB() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connection error: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("DB ping error: %s", err.Error())
	}

	return db
}