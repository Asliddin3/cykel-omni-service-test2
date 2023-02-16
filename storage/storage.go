package storage

import (
	"github.com/Asliddin3/cykel-omni/storage/postgres"
	"github.com/Asliddin3/cykel-omni/storage/repo"
	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Lock() repo.LockStorageI
}

type storagePg struct {
	db       *sqlx.DB
	lockRepo repo.LockStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		lockRepo: postgres.NewLockRepo(db),
	}
}
func (s storagePg) Lock() repo.LockStorageI {
	return s.lockRepo
}
