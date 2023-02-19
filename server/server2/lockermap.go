package server

import (
	"net"
	"sync"

	methods "github.com/Asliddin3/cykel-omni/server/server2/handleMethods/methods"
	repo "github.com/Asliddin3/cykel-omni/server/server2/handleMethods/repo"
)

//LockersMap store all lockers connection
type LockersMap struct {
	Lockers map[int64]repo.LockerMethods
	Mx      sync.RWMutex
}

//NewLocker declared struct with implemented methods and connection
func NewLocker(conn net.Conn) repo.LockerMethods {
	return &methods.Locker{
		UnlockCh:   make(chan string),
		LockerConn: conn,
	}
}

//CheckLockerConn checks locker client connection to server
func (c *LockersMap) CheckLockerConn(imie int64) repo.LockerMethods {
	c.Mx.RLock()
	val, ok := c.Lockers[imie]
	c.Mx.RUnlock()
	if ok {
		return val
	}
	return nil
}

//AddLocker this func add locker struct with methods to map
func (c *LockersMap) AddLocker(imie int64, conn net.Conn) {
	defer c.Mx.Unlock()
	c.Mx.Lock()
	if _, ok := c.Lockers[imie]; ok {
		return
	}
	c.Lockers[imie] = NewLocker(conn)
	return
}

//RemoveConnection remove locker connection from map
func (c *LockersMap) RemoveConnection(imie int64) {
	c.Mx.Lock()
	delete(c.Lockers, imie)
	c.Mx.Unlock()
}
