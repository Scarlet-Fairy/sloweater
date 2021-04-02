package transport

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pb"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedSchedulerServer
	scheduleImageBuild              grpctransport.Handler
	getImageBuildStatus             grpctransport.Handler
	getScheduledImageBuildWorkloads grpctransport.Handler
}

func NewGRPCServer(enpoints service.SchedulerEndpoint, logger log.Logger) pb.SchedulerServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		scheduleImageBuild: grpctransport.NewServer(
			enpoints.ScheduleImageBuildEndpoint,
			decodeScheduleImageBuildRequest,
			encodeScheduleImageBuildResponse,
			options...,
		),
		getImageBuildStatus: grpctransport.NewServer(
			enpoints.GetImageBuildStatusEndpoint,
			decodeGetImageBuildStatusRequest,
			encodeGetImageBuildStatusResponse,
			options...,
		),
		getScheduledImageBuildWorkloads: grpctransport.NewServer(
			enpoints.GetScheduledImageBuildWorkloads,
			decodeGetScheduledImageBuildWorkloadsRequest,
			encodeGetScheduledImageBuildWorkloadsResponse,
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

func (g grpcServer) GetImageBuildStatus(ctx context.Context, request *pb.GetImageBuildStatusRequest) (*pb.GetImageBuildStatusResponse, error) {
	_, resp, err := g.getImageBuildStatus.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetImageBuildStatusResponse), nil
}

func (g grpcServer) GetScheduledImageBuildWorkloads(ctx context.Context, request *pb.GetScheduledImageBuildWorkloadsRequest) (*pb.GetScheduledImageBuildWorkloadsResponse, error) {
	_, resp, err := g.getScheduledImageBuildWorkloads.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetScheduledImageBuildWorkloadsResponse), nil
}

func (g grpcServer) ScheduleWorkload(ctx context.Context, request *pb.ScheduleWorkloadRequest) (*pb.ScheduleWorkloadResponse, error) {
	panic("implement me")
}

func (g grpcServer) GetWorkloadStatus(ctx context.Context, request *pb.GetWorkloadStatusRequest) (*pb.GetWorkloadStatusResponse, error) {
	panic("implement me")
}

func decodeScheduleImageBuildRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ScheduleImageBuildRequest)
	return service.ScheduleImageBuildRequest{
		GitRepoUrl: req.GitRepoUrl,
	}, nil
}

func encodeScheduleImageBuildResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(service.ScheduleImageBuildResponse)
	return &pb.ScheduleImageBuildResponse{
		JobId: res.JobId,
	}, nil
}

func decodeGetImageBuildStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetImageBuildStatusRequest)
	return service.GetImageBuildStatusRequest{
		JobId: req.JobId,
	}, nil
}

func encodeGetImageBuildStatusResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(service.GetImageBuildStatusResponse)

	grpcRes := &pb.GetImageBuildStatusResponse{
		Status: pb.GetImageBuildStatusResponse_ImageBuildStatus(res.Status),
		Steps:  []*pb.GetImageBuildStatusResponse_BuildStepStatus{},
	}

	for _, step := range res.Steps {
		grpcRes.Steps = append(grpcRes.Steps, &pb.GetImageBuildStatusResponse_BuildStepStatus{
			Step:  pb.GetImageBuildStatusResponse_BuildStepStatus_ImageBuildSteps(step.Step),
			Error: &step.Error,
		})
	}

	return grpcRes, nil
}

func decodeGetScheduledImageBuildWorkloadsRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return service.GetScheduledImageBuildWorkloadsRequest{}, nil
}

func encodeGetScheduledImageBuildWorkloadsResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(service.GetScheduledImageBuildWorkloadsResponse)

	return &pb.GetScheduledImageBuildWorkloadsResponse{
		Jobs: res.JobIds,
	}, nil
}
