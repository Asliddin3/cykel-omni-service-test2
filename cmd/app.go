package main

import (
	"fmt"
	"net"
	"os"

	config "github.com/Asliddin3/cykel-omni/config"
	"github.com/Asliddin3/cykel-omni/pkg/logger"
	server "github.com/Asliddin3/cykel-omni/server"
	grpcClient "github.com/Asliddin3/cykel-omni/service/grpc_client"
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
	adminClient, err := grpcClient.New(cfg)
	if err != nil {
		fmt.Println("error connection to admin service", err)
		return
	}
	fmt.Println("Listening on " + cfg.LockHost + ":" + cfg.LockPort)
	tcpChannel := make(chan struct{})
	go server.ListenTCP(l, adminClient, tcpChannel)
	<-tcpChannel
	defer l.Close()
}
