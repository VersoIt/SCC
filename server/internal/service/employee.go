package service

import (
	"serverClientClient/server/internal/model"
	"serverClientClient/server/internal/repository"
)

type EmployeeService struct {
	repo repository.Employee
}

func NewEmployeeService(repo repository.Employee) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetById(id int) (model.Employee, error) {
	return s.repo.GetById(id)
}

func (s *EmployeeService) InitDB(count int) (bool, error) {
	return s.repo.Init(count)
}
