package service

import (
	"serverClientClient/internal/model"
	"serverClientClient/internal/repository"
)

type SocketService struct {
	repo repository.Socket
}

func NewSocketService(repo repository.Socket) *SocketService {
	return &SocketService{repo: repo}
}

func (s *SocketService) SaveToDB(data model.SocketData) error {
	return s.repo.SaveToDB(data)
}
