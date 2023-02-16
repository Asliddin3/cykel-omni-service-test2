package grpcClient

import (
	config "github.com/Asliddin3/cykel-omni/config"
)

//GrpcClientI ...
type ServiceManager struct {
	conf config.Config
}

func New(cnfg config.Config) (*ServiceManager, error) {
	// connReview, err := grpc.Dial(
	// 	fmt.Sprintf("%s:%d", cnfg.ReviewServiceHost, cnfg.ReviewServicePort),
	// 	grpc.WithInsecure(),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("error while dial product service: host: %s and port: %d",
	// 		cnfg.ReviewServiceHost, cnfg.ReviewServicePort)
	// }

	serviceManager := &ServiceManager{
		conf: cnfg,
		// reviewService:   reviewPB.NewReviewServiceClient(connReview),
	}

	return serviceManager, nil
}


// func (s *ServiceManager) ReviewService() reviewPB.ReviewServiceClient {
// 	return s.reviewService
// }
