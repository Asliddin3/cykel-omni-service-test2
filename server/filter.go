package server

import (
	"fmt"
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

func giveResponse(reqArr []string, connectorToGrpc *ConnectLockerToGrpc) (string, error) {
	lockIMEI := reqArr[2]
	timeFormat := reqArr[3]
	lockCommand := reqArr[4]
	fmt.Println("command ", lockCommand,
		" lockIMEI ", lockIMEI)
	imei, err := strconv.Atoi(lockIMEI)
	if err != nil {
		return "", err
	}
	switch lockCommand {
	case "Q0":
		connectorToGrpc.addLocker(int64(imei))
		return checkInCommand(reqArr)
	case "H0":
		return heartBeatCommand(reqArr)
	case "L0":
		resArr := prepareResponse(lockIMEI, timeFormat)
		resArr, err := responseUnlockCommand(reqArr, resArr, connectorToGrpc, int64(imei))
		if err != nil {
			return "", err
		}
		return joinCommand(resArr)
	case "L1":
		resArr := prepareResponse(lockIMEI, timeFormat)
		resArr, err := responseLockCommand(reqArr, resArr)
		if err != nil {
			return "", err
		}
		return joinCommand(resArr)
	case "D0":
		resArr := prepareResponse(lockIMEI, timeFormat)
		resArr, err := responseGetLocation(reqArr, resArr)
		if err != nil {
			return "", err
		}
		return joinCommand(resArr)
	case "D1":
		return responseSetLocationInterval(reqArr)
	default:
		fmt.Println("I don't know this command")
	}
	return "", nil
}

func responseSetLocationInterval(reqArr []string) (string, error) {
	fmt.Println("get response for set interval, get location set ", reqArr[5], " seconds")
	return "", nil
}

func responseGetLocation(reqArr, resArr []string) ([]string, error) {
	locationInden := reqArr[5]
	timeInHourMinSec := reqArr[6]
	locStatus := reqArr[7]
	latitude := reqArr[8]
	latitudeHem := reqArr[9]
	longitude := reqArr[10]
	longitudeHem := reqArr[11]
	sateCount := reqArr[12]
	positionAccuracy := reqArr[13]
	timeInDayMonthYear := reqArr[14]
	altitude := reqArr[15]
	heightUnit := reqArr[16]
	modeIndication := reqArr[17]
	fmt.Println("get location info from locker ", reqArr[2],
		"\n	Location identification ", locationInden,
		"\n UTC time, hhmmss ", timeInHourMinSec,
		"\n Location status ", locStatus,
		"\n Latitude ddmm.mmmm ", latitude,
		"\n Latitude hemisphere ", latitudeHem,
		"\n longitude dddmm.mmmm ", longitude,
		"\n Longitude hemisphere ", longitudeHem,
		"\n Number of satellites searched ", sateCount,
		"\n HDOP ", positionAccuracy,
		"\n UTC date, ddmmyy ", timeInDayMonthYear,
		"\n Altitude ", altitude,
		"\n Height unit M ", heightUnit,
		"\n Mode indication ", modeIndication)
	resArr = append(resArr, "Re", "D0")
	return resArr, nil
}

func responseLockCommand(reqArr, resArr []string) ([]string, error) {
	userID := reqArr[5]
	lockedTime := reqArr[6]
	cycleTime := reqArr[7]
	fmt.Println("locker locked by user ", userID,
		" locked time ", lockedTime, " cycle time ", cycleTime)
	resArr = append(resArr, "Re", "L1")
	return resArr, nil
}

func joinCommand(resArr []string) (string, error) {
	response := strings.Join(resArr, ",")
	response = response + "#\n"
	return response, nil
}
func responseUnlockCommand(reqArr []string, resArr []string, connector *ConnectLockerToGrpc, imie int64) ([]string, error) {

	fmt.Println("get response for unlock command with unlock status",
		reqArr[5], " for user ", reqArr[6], " time ", reqArr[7])
	command := reqArr[2]
	grpcChannel, err := connector.CheckLastCall(imie, command)
	if err != nil {
		return nil, err
	}
	grpcChannel <- strings.Join(reqArr, ",")
	resArr = append(resArr, "Re", "L0")
	return resArr, nil
}

func heartBeatCommand(reqArr []string) (string, error) {
	fmt.Println("get heartBeat command with lock status", reqArr[5],
		" voltage ", reqArr[6], " signal value ", reqArr[7])
	return "", nil
}

func checkInCommand(reqArr []string) (string, error) {
	voltage := reqArr[5]
	fmt.Println("get check-In command with voltage ", voltage, " from ", reqArr[2])
	return "", nil
}
