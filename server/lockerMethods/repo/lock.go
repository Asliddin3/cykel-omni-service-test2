package repo

import (
	pb "github.com/Asliddin3/cykel-omni/genproto/lock"
)

type LockerMethods interface {
	UnlockLocker(*pb.UnlockRequest) (*pb.UnlockResponse, error)
	SendUnlockResponse(string)
}
