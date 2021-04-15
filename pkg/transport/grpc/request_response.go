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
		JobName:   *res.JobName,
		ImageName: *res.ImageName,
	}, nil
}
