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
	LocationCh chan string
}

func (l *Locker) SendUnlockResponse(response string) {
	l.UnlockCh <- response
}
func (l *Locker) SendLocationResponse(response string) {
	l.LocationCh <- response
}

func (l *Locker) UnlockLocker(req *pb.UnlockRequest) (*pb.UnlockResponse, error) {
	imei := strconv.Itoa(int(req.IMEI))
	userID := strconv.Itoa(int(req.UserID))
	var resetTime string
	if req.ResetTime == false {
		resetTime = "1"
	} else {
		resetTime = "0"
	}
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
	// lockerCommand = strings.Replace(lockerCommand, "#", "", 1)
	// lockerCommand = strings.Replace(lockerCommand, "\n", "", 1)

	responseArr := strings.Split(lockerCommand, ",")
	var unlockResult bool
	if responseArr[5] == "0" {
		unlockResult = true
	} else {
		unlockResult = false
	}
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting unlock result to int %v", err)
	}

	userIDStr := responseArr[6]
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting userID to int %v", err)
	}
	unlockStr := strings.TrimFunc(responseArr[7], func(r rune) bool {
		if unicode.IsDigit(r) {
			return false
		}
		return true
	})
	fmt.Println("trimed unlock str ", unlockStr)
	unlockedTime, err := strconv.Atoi(unlockStr)
	if err != nil {
		return &pb.UnlockResponse{}, fmt.Errorf("error converting unlocked time to int %v", err)
	}

	unlockResp := &pb.UnlockResponse{
		UnlockResult: unlockResult,
		UserID:       userIDStr,
		UnlockedTime: int64(unlockedTime),
	}
	return unlockResp, nil
}

func (l *Locker) GetLockerLocation(req *pb.LocationRequest) (*pb.LocationResponse, error) {
	imei := strconv.Itoa(int(req.IMEI))
	reqArr := prepareRequest(imei)
	reqArr = append(reqArr, "D0")
	reqStr := strings.Join(reqArr, ",")
	conn := *l.LockerConn
	_, err := conn.Write([]byte(AddByte([]byte(reqStr))))
	if err != nil {
		return nil, fmt.Errorf("error writing to client connection %v", err)
	}
	response := <-l.UnlockCh
	response = strings.TrimRight(response, "#\n")
	responseArr := strings.Split(response, ",")
	resToGrpc := &pb.LocationResponse{}
	if responseArr[5] == "0" {
		resToGrpc.Tracking = true
	} else {
		resToGrpc.Tracking = false
	}
	resToGrpc.UTCtime = responseArr[6]
	if responseArr[7] != "A" {
		resToGrpc.ValidLocation = false
		resToGrpc.UTCdate = responseArr[14]
	}
	resToGrpc.ValidLocation = true
	latitude, err := strconv.ParseFloat(responseArr[8], 32)
	if err != nil {
		return nil, fmt.Errorf("error converting latitude %v", err)
	}
	resToGrpc.Latitude = float32(latitude)
	if responseArr[9] == "N" {
		resToGrpc.IsNorth = true
	} else {
		resToGrpc.IsNorth = false
	}
	longitude, err := strconv.ParseFloat(responseArr[10], 32)
	if err != nil {
		return nil, fmt.Errorf("error converting longitude %v", err)
	}
	resToGrpc.Longitude = float32(longitude)
	if responseArr[11] == "E" {
		resToGrpc.IsEast = true
	} else {
		resToGrpc.IsEast = false
	}
	sateCount, err := strconv.Atoi(responseArr[12])
	if err != nil {
		return nil, fmt.Errorf("error converting sate count to int %v", err)
	}
	resToGrpc.CountSate = int64(sateCount)
	accuracy, err := strconv.ParseFloat(responseArr[13], 32)
	if err != nil {
		return nil, fmt.Errorf("error converting accuracy to float %v", err)
	}
	resToGrpc.PositionAccuracy = float32(accuracy)
	resToGrpc.UTCdate = responseArr[14]
	altitude, err := strconv.Atoi(responseArr[15])
	if err != nil {
		return nil, fmt.Errorf("error converting altitude to int %v", err)
	}
	resToGrpc.Altitude = int64(altitude)
	resToGrpc.HeightUnit = responseArr[16]
	resToGrpc.ModeIndicatino = responseArr[17]
	return resToGrpc, nil
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
