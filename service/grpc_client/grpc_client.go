package grpcClient

import (
	"fmt"

	config "github.com/Asliddin3/cykel-omni/config"
	adminPB "github.com/Asliddin3/cykel-omni/genproto/admin"
	"google.golang.org/grpc"
)

//ServiceManager ...
type ServiceManager struct {
	conf         config.Config
	adminService adminPB.AdminServiceClient
}

//New this func create service manager
func New(cnfg config.Config) (*ServiceManager, error) {
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cnfg.AdminServiceHost, cnfg.AdminServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial admin service: host: %s and port: %s",
			cnfg.AdminServiceHost, cnfg.AdminServicePort)
	}

	serviceManager := &ServiceManager{
		conf:         cnfg,
		adminService: adminPB.NewAdminServiceClient(connReview),
	}

	return serviceManager, nil
}

//AdminService this func return admin service client
func (s *ServiceManager) AdminService() adminPB.AdminServiceClient {
	return s.adminService
}
