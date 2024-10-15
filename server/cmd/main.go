package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"serverClientClient"
	"serverClientClient/internal/handler"
	"serverClientClient/internal/repository"
	"serverClientClient/internal/service"
	"serverClientClient/pkg/database"
	"strconv"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("error loading .env file")
		return
	}
	gin.SetMode(os.Getenv("GIN_MODE"))

	db, err := database.NewPostgresDB(database.DBConfig{Host: os.Getenv("DB_HOST"), Port: os.Getenv("DB_PORT"), User: os.Getenv("DB_USER"), Password: os.Getenv("DB_PASSWORD"), DBName: os.Getenv("DB_NAME"), SSLMode: os.Getenv("DB_SSLMODE")})
	if err != nil {
		logrus.Errorf("error connecting to the database: %v", err)
		return
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			logrus.Errorf("error closing db connection: %v", err)
		}
	}(db)

	migrator := database.NewMigrator(db)
	if err = migrator.Migrate(database.PostgresDialect, "schema"); err != nil {
		logrus.Error(err)
		return
	}

	srv := server.HttpServer{}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	countToInit, err := strconv.Atoi(os.Getenv("SERVER_EMP_COUNT"))
	if err != nil {
		logrus.Error(err)
		return
	}

	initialized, err := services.Employee.InitDB(countToInit)
	if err != nil {
		logrus.Error(err)
		return
	}

	if initialized {
		logrus.Infof("employees initialized with %d lines", countToInit)
	}

	go func() {
		logrus.Info("server started!")
		err = srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRoutes())
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Info("server stopped")
		} else if err != nil {
			logrus.Error(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("error while stopping server: %s", err)
		return
	}
}
