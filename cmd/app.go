package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	config "github.com/Asliddin3/cykel-omni/config"
	pb "github.com/Asliddin3/cykel-omni/genproto/lock"
	"github.com/Asliddin3/cykel-omni/pkg/logger"
	server "github.com/Asliddin3/cykel-omni/server/server2"
	"github.com/Asliddin3/cykel-omni/server/server2/handleMethods/repo"
	"github.com/Asliddin3/cykel-omni/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "")
	defer logger.Cleanup(log)
	l, err := net.Listen("tcp", cfg.LockHost+":"+cfg.LockPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	lockersMap := &server.LockersMap{
		Lockers: make(map[int64]repo.LockerMethods),
		Mx:      sync.RWMutex{},
	}
	// Close the listener when the application closes.
	fmt.Println("Listening on " + cfg.LockHost + ":" + cfg.LockPort)
	tcpChannel := make(chan struct{})
	go server.ListenTCP(l, lockersMap, tcpChannel)
	lockService := service.NewLockService(nil, lockersMap, log)
	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterLockServiceServer(s, lockService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	<-tcpChannel
	defer l.Close()

}
