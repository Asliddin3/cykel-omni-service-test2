package methods

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"

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
	unlockReqArr := prepareRequest(imei)
	timeInSecond := strconv.Itoa(int(getTimeInSecond()))
	unlockReqArr = append(unlockReqArr, "L0", resetTime, userID, timeInSecond)
	unlockReqByteArr := []byte(strings.Join(unlockReqArr, ","))
	val := *l.LockerConn
	_, err := val.Write(AddByte(unlockReqByteArr))
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error while writing to locker connection: %v", err)
	}
	lockerCommand := <-l.UnlockCh
	lockerCommand = strings.Replace(lockerCommand, "#", "", 1)
	lockerCommand = strings.Replace(lockerCommand, "\n", "", 1)
	fmt.Println("after", lockerCommand)
	responseArr := strings.Split(lockerCommand, ",")
	fmt.Println("data from channel arr", responseArr)
	unlockResult := responseArr[5]
	fmt.Println("converted unlock data from arr to int ", responseArr[5], "--->", unlockResult)
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting unlock result to int %v", err)
	}

	userIDStr := responseArr[6]
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting userID to int %v", err)
	}
	// unlockedTime, err := strconv.Atoi(responseArr[7])
	// if err != nil {
	// 	return &pb.UnlockResponse{}, fmt.Errorf("error converting time to int %v", err)
	// }
	unlockStr := strings.TrimFunc(responseArr[7], func(r rune) bool {
		if unicode.IsDigit(r) {
			return false
		}
		return true
	})
	unlockedTime, err := strconv.Atoi(unlockStr)
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting unlocked time to int %v", err)
	}

	unlockResp := &pb.UnlockResponse{
		UnlockResult: unlockResult,
		UserID:       userIDStr,
		UnlockedTime: string(unlockedTime),
	}
	return unlockResp, nil
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
		fmt.Println("error getting location for time in methods", err)
		return ""
	}
	// timeStr := time.Now().In(loc).Format("20060102150405")
	// timeStr = strings.TrimPrefix(timeStr, "20")

	timeStr := time.Now().In(loc).Format("20060102150405")
	timeStr = "20" + timeStr
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
func getTimeInSecond() int64 {
	loc, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		fmt.Println("error loading loc time for second")
		return 0
	}
	timeInt := time.Now().In(loc).Unix()
	return timeInt
}
