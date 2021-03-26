package transport

import (
	"github.com/Scarlet-Fairy/sloweater/pb"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	scheduleImageBuild              grpctransport.Handler
	getImageBuildStatus             grpctransport.Handler
	getSchedulesImageBuildWorkloads grpctransport.Handler
}

func NewGRPCServer(enpoints service.SchedulerEndpoint, logger log.Logger) pb.SchedulerServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		scheduleImageBuild:              grpctransport.NewServer(),
		getImageBuildStatus:             grpctransport.NewServer(),
		getSchedulesImageBuildWorkloads: grpctransport.NewServer(),
	}
}
