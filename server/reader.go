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
		giveResponse(reqArr, res, lockers)
	}
	ch <- struct{}{}
}

func giveResponse(reqArr []string, reqStr string, lockers *LockersMap) {
	lockIMEI := reqArr[2]
	// timeFormat := reqArr[3]
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
		//send locker check-in to server
		return
	case "H0":
		//send locker heartbeat to sever
		return
	case "L0":
		locker := lockers.Lockers[int64(imei)]
		locker.SendUnlockResponse(reqStr)
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
