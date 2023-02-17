package main

import (
	"fmt"
	"net"
	"os"

	config "github.com/Asliddin3/cykel-omni/config"
	"github.com/Asliddin3/cykel-omni/pkg/logger"
	"github.com/Asliddin3/cykel-omni/server"
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
	commandsMap := &server.ServerCommand{}
	// Close the listener when the application closes.
	fmt.Println("Listening on " + cfg.LockHost + ":" + cfg.LockPort)
	tcpChannel := make(chan struct{})
	go server.ListenTCP(l, commandsMap, tcpChannel)
	<-tcpChannel
	defer l.Close()

}
