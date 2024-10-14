package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"serverClientClient/server/internal"
)

func main() {
	db, err := sqlx.Open("postgres", "user=postgres dbname=db sslmode=disable password=postgres")
	if err != nil {
		logrus.Fatalf("open db err: %v", err)
	}

	if err = internal.InitDB(db); err != nil {
		logrus.Info("db already initialized")
	}
}
