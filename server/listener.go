package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
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
			buf := make([]byte, 1024)
			conn.Read(buf)
			conn.Close()
			cancel()
			return
		}
		go handleRequest(conn, admin, cancel)
	}
	ch <- struct{}{}
}

func handleRequest(conn net.Conn, adminStream pbAdmin.AdminService_LockerStreamingClient, cancel context.CancelFunc) {
	catchError := make(chan error)
	// serverError := make(chan error)
	// var lockerMutex sync.Mutex
	ctx, cancelSubFunc := context.WithCancel(context.Background())
	go recvMessage(ctx, adminStream, conn, catchError)
	go sendMessage(ctx, adminStream, conn, catchError)
	catcherCh := make(chan error)
	// go catchStreamError(clientError, serverError, adminStream, cancel, catcherCh)
	err := <-catcherCh
	fmt.Println(err)
	cancel()
	cancelSubFunc()
	fmt.Println("gotten from catcher channel ", err)
	err = adminStream.CloseSend()
	fmt.Println("gotten from catcher channel ", err)
	err = <-catchError
	fmt.Println(err)

	if err != nil {
		fmt.Println("error while sending close send ", err)
		return
	}
	fmt.Println("stream send closed successfully")

}

func recvMessage(ctx context.Context, recvStream pbAdmin.AdminService_LockerStreamingClient, conn net.Conn, catchError chan error) {
	defer func() {
		conn.Close()
	}()
	for {
		err := recvStream.Context().Err()
		if err != nil {
			fmt.Println("getting error from context ", err)
			catchError <- fmt.Errorf("server error %v", err)
			return
		}
		message, err := recvStream.Recv()
		if err == io.EOF {
			fmt.Println("no more data in stream recv")
			catchError <- fmt.Errorf("server closed sending message")
			return
		} else if err != nil {
			fmt.Println("catched error while recv from stream conn ", err)
			catchError <- fmt.Errorf("error while recovering message from stream %v", err)
			return
		}
		fmt.Println("stream ----->", message.GetAdminMessage())
		// if message.GetAdminMessage() == "" {
		// 	continue
		// }
		fmt.Println("message before writing locker conn ", message.GetAdminMessage())
		_, err = conn.Write(AddByte([]byte(message.GetAdminMessage())))
		if err != nil {
			fmt.Println("chatched error while writing to locker conn ", err)
			catchError <- fmt.Errorf("error while writing to locker connection %v", err)
			return
		}
		fmt.Println("command written to locker conn successfully", message.AdminMessage)
	}
}

func sendMessage(ctx context.Context, sendStream pbAdmin.AdminService_LockerStreamingClient, conn net.Conn, catchError chan error) {
	var lockerIMEI int
	defer sendStream.CloseSend()
	rdr := bufio.NewReader(conn)
	for {
		buf, err := rdr.ReadString('\n')
		fmt.Println("readline result buffer ", buf)
		if err == io.EOF {
			fmt.Println("no more data")
			catchError <- fmt.Errorf("unexpected closed connection from locker")
			return
		} else if err != nil {
			fmt.Println("error while reading from locker conn ", err)
			catchError <- fmt.Errorf("gotten error from locker conn %v", err)
			return
		}
		if buf == "" {
			fmt.Println("gotten empty buffer from locker")
			continue
		}
		if lockerIMEI == 0 {
			lockerIMEI, err = strconv.Atoi(strings.Split(buf, ",")[2])
			if err != nil {
				catchError <- fmt.Errorf("error converting locker imei to int %v", err)
				return
			}
		}
		buf = strings.Replace(buf, "#\n", "", 1)
		err = sendStream.Context().Err()
		if err != nil {
			fmt.Println("getting error from context ", err)
			catchError <- fmt.Errorf("server error %v", err)
			return
		}
		err = sendStream.Send(&pbAdmin.LockerRequest{
			LockerIMEI:    int64(lockerIMEI),
			LockerMessage: buf,
		})
		if err != nil {
			catchError <- fmt.Errorf("error while sending locker request %v", err)
			return
		}
		fmt.Println("sended message to stream ", buf)
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
