package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
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
		ctx, cancel := context.WithCancel(context.Background())
		admin, err := adminClient.AdminService().LockerStreaming(ctx)
		if err != nil {
			fmt.Println("error connection to admin service ", err)
			cancel()
			return
		}
		go handleRequest(conn, admin, cancel)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn, adminStream pbAdmin.AdminService_LockerStreamingClient, cancel context.CancelFunc) {
	clientError := make(chan error)
	serverError := make(chan error)
	// var lockerMutex sync.Mutex

	go recvMessage(adminStream, conn, clientError, serverError)
	go sendMessage(adminStream, conn, clientError, serverError)
	catcherCh := make(chan error)
	go catchStreamError(clientError, serverError, adminStream, cancel, catcherCh)
	err := <-catcherCh
	fmt.Println("gotten from catcher channel ", err)
	err = adminStream.CloseSend()
	if err != nil {
		fmt.Println("error while sending close send ", err)
		return
	}
	fmt.Println("stream send closed successfully")

}

func catchStreamError(clientError chan error, serverError chan error, stream pbAdmin.AdminService_LockerStreamingClient, cancel context.CancelFunc, catcherCh chan error) {
	for {
		select {
		case err := <-clientError:
			catcherCh <- fmt.Errorf("catch client error %v", err)
			cancel()
		case err := <-serverError:
			catcherCh <- fmt.Errorf("catch server error %v", err)
			cancel()
		}
	}
}

func recvMessage(recvStream pbAdmin.AdminService_LockerStreamingClient, conn net.Conn, clientError chan error, serverError chan error) {
	for {
		message, err := recvStream.Recv()
		if err == io.EOF {
			fmt.Println("no more data in stream recv")
			continue
		} else if err != nil {
			serverError <- fmt.Errorf("error while recovering message from stream %v", err)
			return
		}
		fmt.Println("gotten message from stream ", message.GetAdminMessage())
		if message.GetAdminMessage() == "" {
			continue
		}
		// lockerMutex.Lock()
		fmt.Println("before writing message to locker conn ", string(AddByte([]byte(message.AdminMessage))))
		// wr := bufio.NewWriter(conn)
		_, err = conn.Write(AddByte([]byte(message.GetAdminMessage())))
		// defer lockerMutex.Unlock()
		if err != nil {
			clientError <- fmt.Errorf("error while writing to locker connection %v", err)
			return
		}
		fmt.Println("command written to locker conn successfully", message.AdminMessage)
	}
}

func sendMessage(sendStream pbAdmin.AdminService_LockerStreamingClient, conn net.Conn, clientError chan error, serverError chan error) {
	var lockerIMEI int
	rdr := bufio.NewReader(conn)
	for {
		// buf := make([]byte, 1024)
		time.Sleep(time.Second * 1)
		// lockerMutex.Lock()
		buf, err := rdr.ReadString('\n')
		// byteSize, err := conn.Read(buf)
		// lockerMutex.Unlock()
		fmt.Println("readline result bufer ", buf)

		if err != nil {
			clientError <- fmt.Errorf("error while reading from locker connection %v", err)
			return
		}
		res := strings.Replace(buf, "#", "", 1)
		// fmt.Println("gotten command ", string(buf), "with size", byteSize)
		// if byteSize == 0 {
		// 	continue
		// }
		// if lockerIMEI == 0 {
		// 	imeiStr := strings.Split(string(buf[:byteSize]), ",")[2]
		// 	lockerIMEI, err = strconv.Atoi(imeiStr)
		// 	fmt.Println("getting locker imei before stream send", lockerIMEI, err)
		// 	if err != nil {
		// 		clientError <- fmt.Errorf("error converting locker imei to int %v", err)
		// 		return
		// 	}

		// }
		// res := string(buf[:byteSize])
		// res = strings.ReplaceAll(res, "#", "!")
		// res = strings.ReplaceAll(res, "\n", "!")
		// res = strings.ReplaceAll(res, `\`, "")
		// res = strings.TrimRight(res, "!!")
		// fmt.Println("after removing ", res)
		err = sendStream.Send(&pbAdmin.LockerRequest{
			LockerIMEI:    int64(lockerIMEI),
			LockerMessage: res,
		})
		if err != nil {
			serverError <- fmt.Errorf("error while sending locker request %v", err)
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
