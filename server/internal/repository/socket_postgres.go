package repository

import "github.com/jmoiron/sqlx"

type SocketRepository struct {
	db *sqlx.DB
}

func NewSocketRepository(db *sqlx.DB) *SocketRepository {
	
}
