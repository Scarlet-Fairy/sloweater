package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type SchedulerEndpoint struct {
	ScheduleImageBuildEndpoint      endpoint.Endpoint
	GetImageBuildStatusEndpoint     endpoint.Endpoint
	GetScheduledImageBuildWorkloads endpoint.Endpoint
}

func NewEndpoints(s Service) SchedulerEndpoint {
	var scheduleImageBuildEndpoint endpoint.Endpoint
	{
		scheduleImageBuildEndpoint = makeScheduleImageBuildEndpoint(s)
	}

	var getImageBuildStatusEndpoint endpoint.Endpoint
	{
		getImageBuildStatusEndpoint = makeGetImageBuildStatusEndpoint(s)
	}

	var getScheduledImageBuildWorkloads endpoint.Endpoint
	{
		getScheduledImageBuildWorkloads = makeGetScheduledImageBuildWorkloadsEndpoint(s)
	}

	return SchedulerEndpoint{
		ScheduleImageBuildEndpoint:      scheduleImageBuildEndpoint,
		GetImageBuildStatusEndpoint:     getImageBuildStatusEndpoint,
		GetScheduledImageBuildWorkloads: getScheduledImageBuildWorkloads,
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = ScheduleImageBuildResponse{}
	_ endpoint.Failer = GetScheduledImageBuildWorkloadsResponse{}
	_ endpoint.Failer = GetScheduledImageBuildWorkloadsResponse{}
)

type ScheduleImageBuildRequest struct {
	WorkloadId string
	GitRepoUrl string
}

type ScheduleImageBuildResponse struct {
	JobId string `json:"job_id"`
	Err   error  `json:"-"`
}

func (r ScheduleImageBuildResponse) Failed() error {
	return r.Err
}

func makeScheduleImageBuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ScheduleImageBuildRequest)
		jobId, err := s.ScheduleImageBuild(ctx, req.WorkloadId, req.GitRepoUrl)

		return ScheduleImageBuildResponse{
			JobId: jobId,
			Err:   err,
		}, nil
	}
}

type GetImageBuildStatusRequest struct {
	JobId string
}

type GetImageBuildStatusResponse struct {
	Status Status                     `json:"status"`
	Steps  []GetImageBuildStatusSteps `json:"steps"`
	Err    error                      `json:"-"`
}

type GetImageBuildStatusSteps struct {
	Step  Step   `json:"step"`
	Error string `json:"error"`
}

func (r GetImageBuildStatusResponse) Failed() error {
	return r.Err
}

func makeGetImageBuildStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return GetImageBuildStatusResponse{}, nil
	}
}

type GetScheduledImageBuildWorkloadsRequest struct {
}

type GetScheduledImageBuildWorkloadsResponse struct {
	JobIds []string `json:"jobs"`
	Err    error    `json:"-"`
}

func (r GetScheduledImageBuildWorkloadsResponse) Failed() error {
	return r.Err
}

func makeGetScheduledImageBuildWorkloadsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		workloads, err := s.GetSchedulesImageBuildWorkloads(ctx)
		return GetScheduledImageBuildWorkloadsResponse{
			JobIds: workloads,
			Err:    err,
		}, nil
	}
}
