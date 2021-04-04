package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pb"
	"github.com/Scarlet-Fairy/sloweater/pkg/endpoint"
)

func decodeScheduleImageBuildRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ScheduleImageBuildRequest)
	return endpoint.ScheduleImageBuildRequest{
		GitRepoUrl: req.GitRepoUrl,
		WorkloadId: req.WorkloadId,
	}, nil
}

func encodeScheduleImageBuildResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(endpoint.ScheduleImageBuildResponse)
	return &pb.ScheduleImageBuildResponse{
		JobId: res.JobId,
	}, nil
}

func decodeGetImageBuildStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetImageBuildStatusRequest)
	return endpoint.GetImageBuildStatusRequest{
		JobId: req.JobId,
	}, nil
}

func encodeGetImageBuildStatusResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(endpoint.GetImageBuildStatusResponse)

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
	return endpoint.GetScheduledImageBuildWorkloadsRequest{}, nil
}

func encodeGetScheduledImageBuildWorkloadsResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(endpoint.GetScheduledImageBuildWorkloadsResponse)

	return &pb.GetScheduledImageBuildWorkloadsResponse{
		Jobs: res.JobIds,
	}, nil
}
