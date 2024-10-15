package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
	DBName   string
}

const PostgresDialect = "postgres"

func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(PostgresDialect, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
