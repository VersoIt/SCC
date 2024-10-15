package repository

import (
	"github.com/jmoiron/sqlx"
	"serverClientClient/server/internal/model"
)

type PostgresRepository struct {
	Employee
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{Employee: NewEmployeePostgres(db)}
}

type Employee interface {
	GetById(id int) (model.Employee, error)
	Init(count int) (bool, error)
}
