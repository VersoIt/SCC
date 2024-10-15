package service

import (
	"serverClientClient/internal/model"
	"serverClientClient/internal/repository"
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

func (s *EmployeeService) GetAll() ([]model.Employee, error) {
	employees, err := s.repo.GetAll()
	if employees == nil {
		employees = []model.Employee{}
	}
	return employees, err
}
