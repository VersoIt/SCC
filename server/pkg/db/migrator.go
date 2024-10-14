package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db *sqlx.DB
}

func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) migrate(dialect, schemaSrc string) error {
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(m.db.DB, schemaSrc); err != nil {
		return err
	}

	return nil
}
