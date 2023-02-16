package main

import (
	"fmt"
	"net"
	"os"

	config "github.com/Asliddin3/cykel-omni/config"
	"github.com/Asliddin3/cykel-omni/pkg/logger"
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
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + cfg.LockHost + ":" + cfg.LockPort)
	
}
