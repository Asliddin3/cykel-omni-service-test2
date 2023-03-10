package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//ListenTCP this func run loop for connection from client
func ListenTCP(l net.Listener, lockers *LockersMap, ch chan struct{}) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("some error", err)
			break
		}
		go handleRequest(conn, lockers)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn, lockers *LockersMap) {
	bufer := make([]byte, 1024)
	sizeByte, err := conn.Read(bufer)
	if err != nil {
		fmt.Println("error reading from connection")
		return
	}
	if err != nil {

		return
	}
	fmt.Println("first bufer after connection ", string(bufer))
	commands := strings.Split(strings.TrimRight(string(bufer[:sizeByte]), "#\n"), "#\n")
	fmt.Println("trimed  command", commands)
	var lockerIMIE int
	for i := 0; i < len(commands); i++ {
		if i == 0 {
			arr := strings.Split(string(commands[i]), ",")
			lockerIMIE, err := strconv.Atoi(arr[2])
			if err != nil {
				fmt.Println("error converting lockerIMIE to int", err)
				return
			}
			lockers.AddLocker(int64(lockerIMIE), &conn)
			fmt.Println("added connection with imei ", lockerIMIE)
		}
		reqArr := strings.Split(string(commands[i]), ",")
		giveResponse(reqArr, string(commands[i]), lockers, conn)

	}

	readCh := make(chan struct{})
	go ReadClientRequests(conn, readCh, lockers)
	<-readCh
	defer func(conn net.Conn, lockers *LockersMap) {
		fmt.Println("connection removed from map")
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
func prepareRequest(lockIMEI string) []string {
	resArr := make([]string, 4)
	resArr[0] = "*CMDS"
	resArr[1] = "OM"
	resArr[2] = lockIMEI
	resArr[3] = getTime()
	return resArr
}

func getTime() string {
	loc, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		fmt.Println("error loading loc time", err)
		return ""
	}
	// timeStr := time.Now().In(loc).Format("20060102150405")
	// timeStr = strings.TrimPrefix(timeStr, "20")
	timeStr := time.Now().In(loc).Format("20060102150405")
	timeStr = string([]rune(timeStr[2:]))
	return timeStr
	// res := lockerServer.AddByte([]byte(fmt.Sprintf("*CMDS,OM,860537062636022,20200318123020,L0,0,0,%s#\n", timeStr)))
}
