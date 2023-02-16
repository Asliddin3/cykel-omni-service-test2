package postgres

import (
	"github.com/jmoiron/sqlx"
)

type lockRepo struct {
	db *sqlx.DB
}

func NewLockRepo(db *sqlx.DB) *lockRepo {
	return &lockRepo{db: db}
}
