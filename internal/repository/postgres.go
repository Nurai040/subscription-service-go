package repository

import (
	"fmt"
	"os"
	"subscriptions-service/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewPostgresDB() *sqlx.DB {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("POSTGRES_HOST"),
		getEnv("POSTGRES_PORT"),
		getEnv("POSTGRES_USER"),
		getEnv("POSTGRES_PASSWORD"),
		getEnv("POSTGRES_DB"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Log.Fatal("DB connection error", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Log.Fatal("DB ping error", zap.Error(err))
	}

	return db
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		logger.Log.Fatal("missing env variable", zap.String("key", key))
	}
	return v
}