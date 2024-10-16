package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	http2 "serverClientClient/internal/handler/http"
	"serverClientClient/internal/repository"
	"serverClientClient/internal/service"
	"serverClientClient/pkg/database"
	"serverClientClient/pkg/server"
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

	tcpServer, err := server.NewTcpServer(os.Getenv("TCP_SERVER_PORT"), func(conn server.ReadWriteConn) {
		input := make([]byte, 1024)
		bytesRead, err := conn.Read(input)
		if err != nil {
			logrus.Error(err)
		}
		fmt.Println(string(input[:bytesRead]))
		_, err = conn.Write([]byte("Ruslan"))
		if err != nil {
			logrus.Error(err)
		}
	})
	httpServer := server.HttpServer{}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := http2.NewHandler(services)

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
		logrus.Info("TCP server started!")
		err := tcpServer.Run()
		if err != nil && errors.Is(err, net.ErrClosed) {
			logrus.Info("TCP server stopped")
		} else if err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		logrus.Info("HTTP server started!")
		err = httpServer.Run(os.Getenv("HTTP_SERVER_PORT"), handlers.InitRoutes())
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logrus.Info("HTTP server stopped")
		} else if err != nil {
			logrus.Error(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = httpServer.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("error while stopping HTTP server: %s", err)
		return
	}
	err = tcpServer.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("error while stopping TCP server: %s", err)
		return
	}
}
