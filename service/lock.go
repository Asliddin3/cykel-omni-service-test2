package service

import (
	l "github.com/Asliddin3/cykel-omni/pkg/logger"
	grpcclient "github.com/Asliddin3/cykel-omni/service/grpc_client"
	storage "github.com/Asliddin3/cykel-omni/storage"
	"github.com/jmoiron/sqlx"
)

type LockService struct {
	storage storage.IStorage
	client  *grpcclient.ServiceManager
	logger  l.Logger
}

func NewLockService(client *grpcclient.ServiceManager, db *sqlx.DB, log l.Logger) *LockService {
	return &LockService{
		storage: storage.NewStoragePg(db),
		client:  client,
		logger:  log,
	}
}
