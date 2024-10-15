package service

import (
	"serverClientClient/internal/model"
	repository "serverClientClient/internal/repository"
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
	GetAll() ([]model.Employee, error)
}
