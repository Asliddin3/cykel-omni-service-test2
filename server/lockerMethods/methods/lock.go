package methods

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	pb "github.com/Asliddin3/cykel-omni/genproto/lock"
)

const (
	timeFormat string = "200318123020"
)

type Locker struct {
	LockerConn *net.Conn
	UnlockCh   chan string
}

func (l *Locker) SendUnlockResponse(response string) {
	l.UnlockCh <- response
}

func (l *Locker) UnlockLocker(req *pb.UnlockRequest) (*pb.UnlockResponse, error) {
	imei := strconv.Itoa(int(req.IMEI))
	userID := strconv.Itoa(int(req.UserID))
	resetTime := strconv.Itoa(int(req.ResetTime))
	unlockReqArr := prepareRequest(imei, timeFormat)
	unlockReqArr = append(unlockReqArr, "L0", resetTime, userID, getTime())
	unlockReqByteArr := []byte(strings.Join(unlockReqArr, ","))
	val := *l.LockerConn
	fmt.Println("sended command to client connection ", string(AddByte(unlockReqByteArr)))
	_, err := val.Write(AddByte(unlockReqByteArr))
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error while writing to locker connection: %v", err)
	}
	lockerCommand := <-l.UnlockCh
	fmt.Println("gotten from locker channel command ", lockerCommand)
	lockerCommand = strings.TrimRight(lockerCommand, "#\n")
	responseArr := strings.Split(lockerCommand, ",")
	fmt.Println("data from channel arr", responseArr)
	unlockResult, err := strconv.Atoi(responseArr[5])
	fmt.Println("converted unlock data from arr to int ", responseArr[5], "--->", unlockResult)
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting unlock result to int %v", err)
	}

	userIDInt, err := strconv.Atoi(responseArr[6])
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting userID to int %v", err)
	}
	// unlockedTime, err := strconv.Atoi(responseArr[7])
	// if err != nil {
	// 	return &pb.UnlockResponse{}, fmt.Errorf("error converting time to int %v", err)
	// }
	unlockedTime := responseArr[7]
	return &pb.UnlockResponse{
		UnlockResult: int32(unlockResult),
		UserID:       int64(userIDInt),
		UnlockedTime: unlockedTime,
	}, nil
}
func prepareRequest(lockIMEI string, timeFormat string) []string {
	resArr := make([]string, 4)
	resArr[0] = "*CMDS"
	resArr[1] = "OM"
	resArr[2] = lockIMEI
	resArr[3] = timeFormat
	return resArr
}

func getTime() string {
	timeStr := time.Now().Format("20060102150405")
	timeStr = strings.TrimPrefix(timeStr, "20")
	return timeStr
	// res := lockerServer.AddByte([]byte(fmt.Sprintf("*CMDS,OM,860537062636022,20200318123020,L0,0,0,%s#\n", timeStr)))
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
