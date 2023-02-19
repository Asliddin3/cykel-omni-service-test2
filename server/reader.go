package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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
		if err != nil {
			fmt.Println("error reading from client connection", err)
			break
		}
		res := strings.TrimRight(string(buf), "#\n")
		reqArr := strings.Split(res, ",")
		giveResponse(reqArr, res, lockers, conn)
	}
	ch <- struct{}{}
}

func giveResponse(reqArr []string, reqStr string, lockers *LockersMap, conn net.Conn) {
	lockIMEI := reqArr[2]
	timeFormat := reqArr[3]
	lockCommand := reqArr[4]
	fmt.Println("command ", lockCommand,
		" lockIMEI ", lockIMEI)
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
		locker.SendUnlockResponse(reqStr)
		resArr := prepareResponse(lockIMEI, timeFormat)
		resArr = append(resArr, "Re", "L0#\n")
		responseStr := strings.Join(resArr, ",")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to unlock ",responseStr)
		if err != nil {
			fmt.Println("error sending return unlock response")
			return
		}

		return

	default:
		fmt.Println("I don't know this command")
	}
}

func joinCommand(resArr []string) (string, error) {
	response := strings.Join(resArr, ",")
	response = response + "#\n"
	return response, nil
}
