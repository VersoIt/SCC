package repository

import (
	"github.com/jmoiron/sqlx"
	"serverClientClient/internal/model"
)

type PostgresRepository struct {
	Employee
	Socket
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{Employee: NewEmployeePostgres(db), Socket: NewSocketPostgres(db)}
}

type Employee interface {
	GetById(id int) (model.Employee, error)
	Init(count int) (bool, error)
	GetAll() ([]model.Employee, error)
}

type Socket interface {
	SaveToDB(data model.SocketData) error
}
