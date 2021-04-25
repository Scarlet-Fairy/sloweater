package endpoint

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type SchedulerEndpoint struct {
	ScheduleImageBuildEndpoint endpoint.Endpoint
	ScheduleWorkloadEndpoint   endpoint.Endpoint
}

func NewEndpoints(s service.Service, logger log.Logger) SchedulerEndpoint {
	var scheduleImageBuildEndpoint endpoint.Endpoint
	{
		scheduleImageBuildEndpoint = makeScheduleImageBuildEndpoint(s)
		scheduleImageBuildEndpoint = LoggingMiddleware(log.With(logger, "method", "ScheduleImageBuild"))(scheduleImageBuildEndpoint)
		scheduleImageBuildEndpoint = UnwrapErrorMiddleware()(scheduleImageBuildEndpoint)
	}

	var scheduleWorkloadEndpoint endpoint.Endpoint
	{
		scheduleWorkloadEndpoint = makeScheduleWorkloadEndpoint(s)
		scheduleWorkloadEndpoint = LoggingMiddleware(log.With(logger, "method", "ScheduleWorkload"))(scheduleWorkloadEndpoint)
		scheduleWorkloadEndpoint = UnwrapErrorMiddleware()(scheduleWorkloadEndpoint)
	}

	return SchedulerEndpoint{
		ScheduleImageBuildEndpoint: scheduleImageBuildEndpoint,
		ScheduleWorkloadEndpoint:   scheduleWorkloadEndpoint,
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = ScheduleImageBuildResponse{}
)

type ScheduleImageBuildRequest struct {
	WorkloadId string
	GitRepoUrl string
}

type ScheduleImageBuildResponse struct {
	JobName   *string `json:"job_name"`
	ImageName *string `json:"image_name"`
	Err       error   `json:"-"`
}

func (r ScheduleImageBuildResponse) Failed() error {
	return r.Err
}

func makeScheduleImageBuildEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ScheduleImageBuildRequest)
		jobName, imageName, err := s.ScheduleImageBuild(ctx, req.WorkloadId, req.GitRepoUrl)

		return ScheduleImageBuildResponse{
			JobName:   jobName,
			ImageName: imageName,
			Err:       err,
		}, nil
	}
}

type ScheduleWorkloadRequest struct {
	WorkloadId string
	Envs       map[string]string
}

type ScheduleWorkloadResponse struct {
	Err error `json:"-"`
}

func (r ScheduleWorkloadResponse) Failed() error {
	return r.Err
}

func makeScheduleWorkloadEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ScheduleWorkloadRequest)
		err := s.ScheduleWorkload(ctx, req.WorkloadId, req.Envs)

		return ScheduleWorkloadResponse{
			Err: err,
		}, nil
	}
}
