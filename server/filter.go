package server

import (
	"fmt"
	"strings"
)

func giveResponse(req string) (string, error) {
	resArr := make([]string, 4)
	resArr[0] = "*CMDS"
	resArr[1] = "OM"
	req = strings.TrimLeft(req, "#\n")
	reqArr := strings.Split(req, ",")
	lockIMEI := reqArr[2]
	timeFormat := reqArr[3]
	lockCommand := reqArr[4]
	resArr[2] = lockIMEI
	resArr[3] = timeFormat
	fmt.Println("command ", lockCommand,
		" lockIMEI ", lockIMEI)

	switch lockCommand {
	case "Q0":
		return checkInCommand(reqArr)
	case "H0":
		return heartBeatCommand(reqArr)
	case "L0":
		resArr, err := responseUnlockCommand(reqArr, resArr)
		if err != nil {
			return "", err
		}
		return sendCommand(resArr)
	case "L1":
		resArr, err := responseLockCommand(reqArr, resArr)
		if err != nil {
			return "", err
		}
		return sendCommand(resArr)
	case "D0":
		resArr, err := responseGetLocation(reqArr, resArr)
		if err != nil {
			return "", err
		}
		return sendCommand(resArr)
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

func sendCommand(resArr []string) (string, error) {
	response := strings.Join(resArr, ",")
	response = response + "#\n"
	return response, nil
}
func responseUnlockCommand(reqArr []string, resArr []string) ([]string, error) {
	fmt.Println("get response for unlock command with unlock status",
		reqArr[5], " for user ", reqArr[6], " time ", reqArr[7])
	reqArr = append(reqArr, "Re", "L0")
	return reqArr, nil
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
