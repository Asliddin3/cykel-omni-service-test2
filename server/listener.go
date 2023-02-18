package server

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//ListenTCP this func run loop for connection from client
func ListenTCP(l net.Listener, commands *ConnectLockerToGrpc, ch chan struct{}) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("some error", err)
			break
		}
		conn.RemoteAddr()
		go handleRequest(conn, commands)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn, commands *ConnectLockerToGrpc) {
	buf := make([]byte, 2048)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	req := strings.TrimLeft(string(buf), "#\n")
	reqArr := strings.Split(req, ",")
	response, err := giveResponse(reqArr)
	lockIMEI := reqArr[2]
	imei, err := strconv.Atoi(lockIMEI)
	if err != nil {
		fmt.Println("error converting lockIMEI to int")
		return
	}
	grpcChannel := commands.GetChannel(int64(imei))
	grpcChannel <- string(buf)

	if err != nil {
		fmt.Println("Error filtering request data", err)
	}
	if response != "" {
		res := AddByte([]byte(response))
		fmt.Println("send message", string(res))
		_, err = conn.Write([]byte(res))
		if err != nil {
			fmt.Println("write error", err)
		}
	} else {
		fmt.Println("nothing to send to ", string(buf))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*12)
	defer cancel()
	err = waitServerCommand(ctx, conn, lockIMEI, commands)
	if err != nil {
		fmt.Println("waiting grpc command error", err)
		recover()
	}
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()
	// waitServerCommand(ctx, conn)

	defer conn.Close()
}

func waitServerCommand(ctx context.Context, conn net.Conn, lockImei string, lockerChannel *ConnectLockerToGrpc) error {
	imei, err := strconv.Atoi(lockImei)
	if err != nil {
		fmt.Println("error converting lockImei to int", err)
		return err
	}
	// temp := &grpcToLocker{}
	commands := lockerChannel.GetCommands(int64(imei))
	for _, command := range commands {
		_, err = conn.Write([]byte(command))
		if err != nil {
			return err
		}
	}
	fmt.Println("server commands sended successfully")
	return nil
}

//AddByte this func will add two 0xFF byte before command
func AddByte(b2 []byte) []byte {
	arrByte := make([]byte, 2)
	arrByte[0] = 0xFF
	arrByte[1] = 0xFF
	arrByte = append(arrByte, b2...)
	return arrByte
}
