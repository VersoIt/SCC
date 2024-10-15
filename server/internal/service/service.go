package service

import (
	"serverClientClient/server/internal/model"
	"serverClientClient/server/internal/repository"
)

type Service struct {
	Employee
}

func NewService(repo *repository.PostgresRepository) *Service {
	return &Service{Employee: NewEmployeeService(repo)}
}

type Employee interface {
	GetById(id int) (model.Employee, error)
	InitDB(count int) (bool, error)
}
