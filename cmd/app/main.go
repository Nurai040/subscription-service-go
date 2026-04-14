package main

import (
	"fmt"
	"os"
	"strings"

	_ "subscriptions-service/docs"

	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/repository"
	"subscriptions-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		logger.Log.Fatal("missing env variable", zap.String("key", key))
	}
	return v
}

func main() {
	logger.Init()
	defer logger.Log.Sync()

	// env
	user := mustGetenv("POSTGRES_USER")
	pass := mustGetenv("POSTGRES_PASSWORD")
	host := mustGetenv("POSTGRES_HOST")
	port := mustGetenv("POSTGRES_PORT")
	dbName := mustGetenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbName,
	)

	// DB init (оставляем без DSN — как у тебя реализовано внутри)
	db := repository.NewPostgresDB()

	// migrations
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		logger.Log.Fatal("migration init error", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal("migration failed", zap.Error(err))
	}

	// service layer
	repo := repository.NewSubscriptionRepo(db)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc)

	logger.Log.Info("starting server")

	r := gin.Default()

	r.POST("/subscriptions", h.Create)
	r.GET("/subscriptions", h.GetAll)
	r.GET("/subscriptions/:id", h.GetByID)
	r.PUT("/subscriptions/:id", h.Update)
	r.DELETE("/subscriptions/:id", h.Delete)
	r.GET("/subscriptions/summary", h.GetSummary)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("server failed", zap.Error(err))
	}
}