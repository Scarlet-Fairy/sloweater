package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pb"
	"github.com/Scarlet-Fairy/sloweater/pkg/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedSchedulerServer
	scheduleImageBuild grpctransport.Handler
	scheduleWorkload   grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.SchedulerEndpoint, logger log.Logger) pb.SchedulerServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		scheduleImageBuild: grpctransport.NewServer(
			endpoints.ScheduleImageBuildEndpoint,
			decodeScheduleImageBuildRequest,
			encodeScheduleImageBuildResponse,
			options...,
		),
		scheduleWorkload: grpctransport.NewServer(
			endpoints.ScheduleWorkloadEndpoint,
			decodeScheduleWorkloadRequest,
			encodeScheduleWorkloadResponse,
			options...,
		),
	}
}

func (g grpcServer) ScheduleImageBuild(ctx context.Context, request *pb.ScheduleImageBuildRequest) (*pb.ScheduleImageBuildResponse, error) {
	_, resp, err := g.scheduleImageBuild.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ScheduleImageBuildResponse), nil

}

func (g grpcServer) ScheduleWorkload(ctx context.Context, request *pb.ScheduleWorkloadRequest) (*pb.ScheduleWorkloadResponse, error) {
	_, resp, err := g.scheduleWorkload.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ScheduleWorkloadResponse), nil
}
