package server

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	pbAdmin "github.com/Asliddin3/cykel-omni/genproto/admin"
	grpcClient "github.com/Asliddin3/cykel-omni/service/grpc_client"
)

//ListenTCP this func run loop for connection from client
func ListenTCP(l net.Listener, adminClient *grpcClient.ServiceManager, ch chan struct{}) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("some error", err)
			break
		}
		admin, err := adminClient.AdminService().LockerStreaming(context.Background())
		if err != nil {
			fmt.Println("error connection to admin service ", err)
			return
		}
		go handleRequest(conn, admin)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn, adminStream pbAdmin.AdminService_LockerStreamingClient) {
	errorCh := make(chan error)
	var lockerMutex sync.Mutex
	go recvMessage(adminStream, conn, errorCh, &lockerMutex)
	go sendMessage(adminStream, conn, errorCh, &lockerMutex)
	err := <-errorCh
	fmt.Println("error in stream process ", err)
	err = adminStream.CloseSend()
	errClosing := <-errorCh
	fmt.Println("error about closing connection ", errClosing)
	if err != nil {
		fmt.Println("error while sending close send ", err)
		return
	}

}

func recvMessage(recvStream pbAdmin.AdminService_LockerStreamingClient, conn net.Conn, errorCh chan error, lockerMutex *sync.Mutex) {
	for {
		message, err := recvStream.Recv()
		if err != nil {
			errorCh <- fmt.Errorf("error while recovering message from stream %v", err)
			return
		}
		fmt.Println("gotten message from stream ", message.AdminMessage)
		lockerMutex.Lock()
		defer lockerMutex.Unlock()
		if message.AdminMessage == "" {
			continue
		}
		_, err = conn.Write(AddByte([]byte(message.AdminMessage)))
		if err != nil {
			errorCh <- fmt.Errorf("error while writing to locker connection %v", err)
			return
		}
	}

}

func sendMessage(sendStream pbAdmin.AdminService_LockerStreamingClient, conn net.Conn, errorCh chan error, lockerMutex *sync.Mutex) {
	var lockerIMEI int
	for {
		buf := make([]byte, 1024)
		lockerMutex.Lock()
		_, err := conn.Read(buf)
		lockerMutex.Unlock()
		if err != nil {
			errorCh <- fmt.Errorf("error while reading from locker connection %v", err)
			return
		}
		if lockerIMEI == 0 {
			imeiStr := strings.Split(string(buf), ",")[2]
			lockerIMEI, err = strconv.Atoi(imeiStr)
			if err != nil {
				errorCh <- fmt.Errorf("error converting locker imei to int %v", err)
				return
			}
		}
		arrRune := []rune(string(buf[:]))
		res := string(arrRune[:len(arrRune)-2])
		err = sendStream.Send(&pbAdmin.LockerRequest{
			LockerIMEI:    int64(lockerIMEI),
			LockerMessage: res,
		})
		if err != nil {
			errorCh <- fmt.Errorf("error while sending locker request %v", err)
			return
		}
		fmt.Println("sended message to stream ", res)
	}
}

//AddByte this func will add two 0xFF byte before command
func AddByte(b2 []byte) []byte {
	arrByte := make([]byte, 2)
	arrByte[0] = 0xFF
	arrByte[1] = 0xFF
	arrByte = append(arrByte, b2...)
	arrByte = append(arrByte, []byte("#\n")...)
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
	timeStr = "20" + timeStr
	return timeStr
	// res := lockerServer.AddByte([]byte(fmt.Sprintf("*CMDS,OM,860537062636022,20200318123020,L0,0,0,%s#\n", timeStr)))
}
