package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/Asliddin3/cykel-omni/genproto/lock"
	l "github.com/Asliddin3/cykel-omni/pkg/logger"
	server "github.com/Asliddin3/cykel-omni/server"
	storage "github.com/Asliddin3/cykel-omni/storage"
	"github.com/jmoiron/sqlx"
)

//LockService this struct store necessary environments for lock service
type LockService struct {
	storage storage.IStorage
	lockers *server.LockersMap
	logger  l.Logger
}

//NewLockService connect to services for grpc methods
func NewLockService(db *sqlx.DB, lockers *server.LockersMap, log l.Logger) *LockService {
	return &LockService{
		storage: storage.NewStoragePg(db),
		lockers: lockers,
		logger:  log,
	}
}

//UnlockLocker connect to locker tcp client connection and send command then get response
func (l *LockService) UnlockLocker(ctx context.Context, req *pb.UnlockRequest) (*pb.UnlockResponse, error) {
	if l.lockers.CheckLockerConn(req.IMEI) == false {
		return &pb.UnlockResponse{}, fmt.Errorf("no such imie locker connected")
	}
	locker, _ := l.lockers.Lockers[req.IMEI]
	unlockRes, err := locker.UnlockLocker(req)
	if err != nil {
		return &pb.UnlockResponse{}, err
	}
	return unlockRes, nil
}

//For server1
//UnlockLocker this func will unlock locker
// func (l *LockService) UnlockLocker(ctx context.Context, req *pb.UnlockRequest) (*pb.UnlockResponse, error) {
// 	request := fmt.Sprintf("*CMDS,OM,%d,20200318123020,L0,%d,%d,%s#\n", req.IMEI, req.ResetTime, req.UserID, getTime())
// 	lockerRequest := make(chan string)
// 	lockerResponse := make(chan string)

// 	ticker := time.NewTicker(30 * time.Second)
// 	err := l.lockConnector.AddCommand(req.IMEI, request, lockerRequest)
// 	if err != nil {
// 		return nil, fmt.Errorf("error adding command to map %w", err)
// 	}
// 	// Creating channel using make
// 	// reqErr := make(chan error)
// 	go func(requestCh <-chan string, responseCh chan<- string) {
// 		for {
// 			select {
// 			case result := <-requestCh:
// 				responseCh <- result
// 				return
// 			case <-ticker.C:
// 				responseCh <- ""
// 				return
// 			}
// 		}
// 	}(lockerRequest, lockerResponse)

// 	// time.Sleep(31 * time.Second)
// 	result := <-lockerResponse
// 	if result == "" {
// 		lockerRequest <- ""
// 	} else {

// 	}
// 	// ticker.Stop()

// 	return &pb.UnlockResponse{}, nil
// }

// *CMDS,OM,860537062636022,20200318123020,L0,0,0,%s#\n
func getTime() string {
	timeStr := time.Now().Format("20060102150405")
	timeStr = strings.TrimPrefix(timeStr, "20")
	return timeStr
	// res := lockerServer.AddByte([]byte(fmt.Sprintf("*CMDS,OM,860537062636022,20200318123020,L0,0,0,%s#\n", timeStr)))
}
