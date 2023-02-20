package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat string = "200318123020"
)

func prepareResponse(lockIMEI string) []string {
	resArr := make([]string, 4)
	resArr[0] = "*CMDS"
	resArr[1] = "OM"
	resArr[2] = lockIMEI
	resArr[3] = getTime()
	return resArr
}

//ReadClientRequests run loop for reading data from client connection
func ReadClientRequests(conn net.Conn, ch chan struct{}, lockers *LockersMap) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading from client connection ", err)
			break
		}
		res := strings.TrimRight(string(buf), "#\n")
		reqArr := strings.Split(res, ",")
		giveResponse(reqArr, res, lockers, conn)
		time.Sleep(time.Second * 1)
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
		responseStr := makeReturn(lockIMEI, "L0")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to unlock ", responseStr)
		if err != nil {
			fmt.Println("error sending return unlock response")
			return
		}
		return
	case "L1":
		// there should be implemetation for lock command
		responseStr := makeReturn(lockIMEI, "L1")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to lock ", responseStr)
		if err != nil {
			fmt.Println("error sending return unlock response")
			return
		}
	case "D0":
		responseStr := makeReturn(lockIMEI, "D0")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to getting position ", responseStr)
		if err != nil {
			fmt.Println("error sending return  getting position response")
			return
		}
	case "W0":
		responseStr := makeReturn(lockIMEI, "W0")
		_, err = conn.Write(AddByte([]byte(responseStr)))
		fmt.Println("sended return to getting position ", responseStr)
		if err != nil {
			fmt.Println("error sending return  getting position response")
			return
		}
	case "U0":
		res := fmt.Sprintf("*CMDS,OM,%s,%s,U0,220,128,32434,A1,h566m", lockIMEI, getTime())
		_, err := conn.Write([]byte(res))
		if err != nil {
			fmt.Println("error writing upgrade command", err)
			return
		}
	default:
		fmt.Println("I don't know this command")
	}
}
func makeReturn(lockIMEI string, command string) string {
	resArr := prepareResponse(lockIMEI)
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
