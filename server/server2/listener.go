package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

//ListenTCP this func run loop for connection from client
func ListenTCP(l net.Listener, lockers *LockersMap, ch chan struct{}) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("some error", err)
			break
		}
		conn.RemoteAddr()
		go handleRequest(conn, lockers)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn, lockers *LockersMap) {
	bufer := make([]byte, 1024)
	_, err := conn.Read(bufer)
	if err != nil {
		fmt.Println("error reading from connection")
		return
	}
	arr := strings.Split(string(bufer), ",")
	lockerIMIE, err := strconv.Atoi(arr[2])
	if err != nil {
		fmt.Println("error converting lockerIMIE to int", err)
		return
	}
	lockers.AddLocker(int64(lockerIMIE), conn)
	readCh := make(chan struct{})
	go ReadClientRequests(conn, readCh,lockers)
	<-readCh
	defer func(conn net.Conn, lockers *LockersMap) {
		lockers.RemoveConnection(int64(lockerIMIE))
		conn.Close()
	}(conn, lockers)
}

//AddByte this func will add two 0xFF byte before command
func AddByte(b2 []byte) []byte {
	arrByte := make([]byte, 2)
	arrByte[0] = 0xFF
	arrByte[1] = 0xFF
	arrByte = append(arrByte, b2...)
	return arrByte
}
