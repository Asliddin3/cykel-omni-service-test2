package server

import (
	"fmt"
	"strings"
	"sync"
)

type Connect struct {
	Commands []string
	Channels []chan string
}
type grpcToLocker struct {
	// unlock   unlockLockRequest
}

// type lockerToGrpc struct {
// }

// type unlockLockRequest struct {
// 	lockerImei int
// 	userID     int
// }

// type unlockLockResponse struct {
// 	status    int8
// 	userID    int
// 	timestamp string
// }

//ConnectLockerToGrpc map for saving commands from server
type ConnectLockerToGrpc struct {
	Connector map[int64]Connect
	mx        sync.RWMutex
}

func (c *ConnectLockerToGrpc) addLocker(imie int64) {
	defer c.mx.Unlock()
	c.mx.Lock()
	if _, ok := c.Connector[imie]; ok {
		return
	}
	c.Connector[imie] = Connect{
		Commands: []string{},
		Channels: []chan string{},
	}
	return
}

func (c *ConnectLockerToGrpc) AddCommand(imei int64, command string, ch chan string) error {
	defer c.mx.Unlock()
	c.mx.Lock()
	if lockConn, ok := c.Connector[imei]; ok {
		lockConn.Channels = append(lockConn.Channels, ch)
		lockConn.Commands = append(lockConn.Commands, command)
		c.Connector[imei] = lockConn
		return nil
	}
	return fmt.Errorf("this %d locker never been connected", imei)
}

func (c *ConnectLockerToGrpc) CheckLastCall(imie int64, command string) (chan string, error) {
	defer c.mx.RUnlock()
	c.mx.RLock()
	lockerConnector := c.Connector[imie]
	for i := len(lockerConnector.Commands) - 1; i >= 0; i++ {
		if strings.Contains(lockerConnector.Commands[i], command) {
			return lockerConnector.Channels[i], nil
		}
	}
	return nil, fmt.Errorf("This command does't exists in map commands array")
}

func (c *ConnectLockerToGrpc) RemoveCall(imie int64) error {
	defer c.mx.Unlock()
	c.mx.Lock()
	if val, ok := c.Connector[imie]; ok {
		val.Channels = val.Channels[:len(val.Channels)-2]
		val.Commands = val.Commands[:len(val.Commands)-2]
		c.Connector[imie] = val
		return nil
	}
	return fmt.Errorf("no such kind imie lock")
}

func (c *ConnectLockerToGrpc) GetChannel(imie int64) chan string {
	c.mx.RLock()
	defer c.mx.RUnlock()
	if val, ok := c.Connector[imie]; ok {
		return val.Channels[len(val.Channels)-1]
	}
	return make(chan string)
}

func (c *ConnectLockerToGrpc) GetCommands(imie int64) []string {
	c.mx.RLock()
	defer c.mx.RUnlock()
	if val, ok := c.Connector[imie]; ok {
		return val.Commands
	}
	return []string{}
}
