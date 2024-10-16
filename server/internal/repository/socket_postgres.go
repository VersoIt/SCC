package repository

import (
	"github.com/jmoiron/sqlx"
	"serverClientClient/internal/model"
)

type SocketPostgres struct {
	db *sqlx.DB
}

func NewSocketPostgres(db *sqlx.DB) *SocketPostgres {
	return &SocketPostgres{db: db}
}

func (p *SocketPostgres) SaveToDB(socket model.SocketData) error {
	_, err := p.db.Exec("INSERT INTO sockets (id, data) VALUES($1, $2) ON CONFLICT (id) DO UPDATE SET id = $1, data = $2", socket.Id, socket.Data)
	if err != nil {
		return err
	}
	return nil
}
