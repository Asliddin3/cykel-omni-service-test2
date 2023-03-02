package grpcClient

import (
	"fmt"

	config "github.com/Asliddin3/cykel-omni/config"
	lockerPB "github.com/Asliddin3/cykel-omni/genproto/locker"
	"google.golang.org/grpc"
)

//ServiceManager ...
type ServiceManager struct {
	conf          config.Config
	lockerService lockerPB.LockerServiceClient
}

//New this func create service manager
func New(cnfg config.Config) (*ServiceManager, error) {
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cnfg.LockerServiceHost, cnfg.LockerServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial admin service: host: %s and port: %s",
			cnfg.LockerServiceHost, cnfg.LockerServicePort)
	}

	serviceManager := &ServiceManager{
		conf:          cnfg,
		lockerService: lockerPB.NewLockerServiceClient(connReview),
	}

	return serviceManager, nil
}

//LockerService this func return admin service client
func (s *ServiceManager) LockerService() lockerPB.LockerServiceClient {
	return s.lockerService
}
