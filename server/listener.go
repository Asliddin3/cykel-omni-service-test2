package server

import (
	"context"
	"fmt"
	"net"
	"sync"
)
//ServerCommand map for saving commands from server
type ServerCommand struct {
	lockerCommands map[int][]string
	mx             sync.RWMutex
}

//AddCommand this command will add command to array of commands in map
func (c *ServerCommand) AddCommand(key int, command string) {
	c.lockerCommands[key] = append(c.lockerCommands[key], command)
}

//Load this command with load val from map by key
func (c *ServerCommand) Load(key int) ([]string, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.lockerCommands[key]
	return val, ok
}

//ListenTCP this func run loop for connection from client
func ListenTCP(l net.Listener, commands *ServerCommand, ch chan struct{}) {
	for {
		conn, err := l.Accept()
		fmt.Println("accepted error", err)
		if err != nil {
			fmt.Println("some error", err)
			break
		}
		go handleRequest(conn)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 2048)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	response, err := giveResponse(string(buf))
	if err != nil {
		fmt.Println("Error filtering request data", err)
	}
	// timeStr := time.Now().Format("20060102150405")
	// timeStr = strings.TrimPrefix(timeStr, "20")
	// res := addByte([]byte(fmt.Sprintf("*CMDS,OM,860537062636022,20200318123020,L0,0,0,%s#\n", timeStr)))
	if response != "" {
		res := addByte([]byte(response))
		fmt.Println("send message", string(res))
		_, err = conn.Write([]byte(res))
		if err != nil {
			fmt.Println("write error", err)
		}
	} else {
		fmt.Println("nothing to send to ", string(buf))
	}
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()
	// waitServerCommand(ctx, conn)

	conn.Close()
}

func waitServerCommand(ctx context.Context, conn net.Conn) {

}

func addByte(b2 []byte) []byte {
	arrByte := make([]byte, 2)
	arrByte[0] = 0xFF
	arrByte[1] = 0xFF
	arrByte = append(arrByte, b2...)
	return arrByte
}
