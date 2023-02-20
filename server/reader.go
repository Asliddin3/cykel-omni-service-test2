package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat string = "200318123020"
)

func prepareResponse(lockIMEI string, timeFormat string) []string {
	resArr := make([]string, 4)
	resArr[0] = "*CMDS"
	resArr[1] = "OM"
	resArr[2] = lockIMEI
	resArr[3] = timeFormat
	return resArr
}

//ReadClientRequests run loop for reading data from client connection
func ReadClientRequests(conn net.Conn, ch chan struct{}, lockers *LockersMap) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)

		if err == io.EOF {
			fmt.Println("error reading from client connection ", err)
			time.Sleep(time.Second * 1)
			continue
		} else if err != nil {
			fmt.Println("error reading from client connection ", err)
			time.Sleep(time.Second * 1)
			continue
		}
		res := strings.TrimRight(string(buf), "#\n")
		reqArr := strings.Split(res, ",")
		giveResponse(reqArr, res, lockers, conn)
	}
	ch <- struct{}{}
}

func giveResponse(reqArr []string, reqStr string, lockers *LockersMap, conn net.Conn) {
	lockIMEI := reqArr[2]
	// timeFormat := reqArr[3]
	lockCommand := reqArr[4]
	fmt.Println("command ", lockCommand,
		" lockIMEI ", lockIMEI)
	fmt.Println("gotten command <----", reqStr)
	imei, err := strconv.Atoi(lockIMEI)
	if err != nil {
		fmt.Println("error converting lock imei to int", err)
		return
	}
	switch lockCommand {
	case "Q0":
		fmt.Println("get check in command")
		//send locker check-in to server
		return
	case "H0":
		fmt.Println("get heartbeat command")
		//send locker heartbeat to sever
		return
	case "L0":
		locker := lockers.Lockers[int64(imei)]
		// if timeFormat != "000000000000" && reqArr[5] != "0" && reqArr[6] != "0" {
		fmt.Println("send locker response to unlock channel", reqStr)
		locker.SendUnlockResponse(reqStr)
		// }
		responseStr := makeReturn(lockIMEI, timeFormat, "L0")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to unlock ", responseStr)
		if err != nil {
			fmt.Println("error sending return unlock response")
			return
		}
		return
	case "L1":
		// there should be implemetation for lock command
		responseStr := makeReturn(lockIMEI, timeFormat, "L1")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to lock ", responseStr)
		if err != nil {
			fmt.Println("error sending return unlock response")
			return
		}
	case "D0":
		responseStr := makeReturn(lockIMEI, timeFormat, "D0")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to getting position ", responseStr)
		if err != nil {
			fmt.Println("error sending return  getting position response")
			return
		}
	case "W0":
		responseStr := makeReturn(lockIMEI, timeFormat, "W0")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to getting position ", responseStr)
		if err != nil {
			fmt.Println("error sending return  getting position response")
			return
		}
	default:
		fmt.Println("I don't know this command")
	}
}
func makeReturn(lockIMEI string, timeFormat string, command string) string {
	resArr := prepareResponse(lockIMEI, timeFormat)
	command = command + "#\n"
	resArr = append(resArr, "Re", command)
	responseStr := strings.Join(resArr, ",")
	return responseStr
}

func joinCommand(resArr []string) (string, error) {
	response := strings.Join(resArr, ",")
	response = response + "#\n"
	return response, nil
}
