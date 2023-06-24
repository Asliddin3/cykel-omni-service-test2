package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	pbLocker "github.com/Asliddin3/cykel-omni/genproto/locker"
	grpcClient "github.com/Asliddin3/cykel-omni/service/grpc_client"
)

type LockerStream struct {
	connection map[int64]net.Conn
	Mx         sync.RWMutex
}

type LockerStreams struct {
	Streams [10]*LockerStream
	Mx      sync.RWMutex
}

//ListenTCP this func run loop for connection from client
func ListenTCP(l net.Listener, adminClient *grpcClient.ServiceManager, ch chan struct{}) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("some error", err)
			break
		}
		ctx, cancel := context.WithCancel(context.Background())
		admin, err := adminClient.LockerService().LockerStreaming(ctx)
		if err != nil {
			fmt.Println("error connection to admin service ", err)
			buf := make([]byte, 1024)
			conn.Read(buf)
			conn.Close()
			cancel()
			return
		}
		go handleRequest(ctx, conn, admin, cancel)
	}
	ch <- struct{}{}
}

func handleRequest(ctx context.Context, conn net.Conn, adminStream pbLocker.LockerService_LockerStreamingClient, cancel context.CancelFunc) {
	// serverError := make(chan error)
	// var lockerMutex sync.Mutex
	// ctx, cancelSubFunc := context.WithCancel(context.Background())
	write := make(chan string)
	read := make(chan string)
	var wg sync.WaitGroup
	wg.Add(4)
	go recvMessage(adminStream, write, &wg)
	go writeMessage(conn, write, &wg)
	go sendMessage(adminStream, read, &wg)
	go readMessage(conn, read, &wg)
	// go catchStreamError(clientError, serverError, adminStream, cancel, catcherCh)
	for {
		select {
		case <-adminStream.Context().Done():
			return
		}
	}
}

func writeMessage(conn net.Conn, write chan string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		conn.Close()
	}()
	for {
		message, ok := <-write
		if ok == false {
			conn.Close()
			return
		}
		fmt.Println("message before writing locker conn ", message)
		_, err := conn.Write(AddByte([]byte(message)))
		if err != nil {
			fmt.Println("got error while writing ", err)
			conn.Close()
			return
		}
	}
}

func recvMessage(recvStream pbLocker.LockerService_LockerStreamingClient, write chan string, wg *sync.WaitGroup) {
	defer func() {
		close(write)
		wg.Done()
	}()
	for {
		err := recvStream.Context().Err()
		if err != nil {
			fmt.Println("getting error from context ", err)
			return
		}
		message, err := recvStream.Recv()
		if err == io.EOF {
			fmt.Println("no more data in stream recv")
			return
		} else if err != nil {
			fmt.Println("catched error while recv from stream conn ", err)
			return
		}
		fmt.Println("stream ----->", message.GetStreamMessage())
		channelMessage := message.GetStreamMessage()
		write <- channelMessage
	}
}

func readMessage(conn net.Conn, read chan string, wg *sync.WaitGroup) {
	rdr := bufio.NewReader(conn)
	defer func() {
		conn.Close()
		wg.Done()
	}()
	for {
		buf, err := rdr.ReadString('\n')
		fmt.Println("readline result buffer ", buf)
		if err == io.EOF {
			fmt.Println("no more data")
			return
		} else if err != nil {
			fmt.Println("error while reading from locker conn ", err)
			return
		}

		buf = strings.Replace(buf, "#\n", "", 1)
		read <- buf
	}
}

func sendMessage(sendStream pbLocker.LockerService_LockerStreamingClient, read chan string, wg *sync.WaitGroup) {
	var lockerIMEI int
	defer func() {
		sendStream.CloseSend()
		wg.Done()
	}()
	buffer, ok := <-read
	if !ok {
		fmt.Println("read from closed channel")
		return
	}
	lockerIMEI, err := strconv.Atoi(strings.Split(buffer, ",")[2])
	if err != nil {
		fmt.Println("error converting imei to int")
		return
	}
	for {
		err := sendStream.Context().Err()
		if err != nil {
			fmt.Println("getting error from context ", err)
			return
		}
		buf, ok := <-read
		if !ok {
			fmt.Println("read from closed channel error ")
			return
		}
		err = sendStream.Send(&pbLocker.LockerRequest{
			LockerIMEI:    int64(lockerIMEI),
			LockerMessage: buf,
		})
		if err != nil {
			fmt.Println("catched error while sending message")
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
