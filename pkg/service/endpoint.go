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
}

type ScheduleImageBuildResponse struct {
	Err error `json:"-"`
}

func (r ScheduleImageBuildResponse) Failed() error {
	return r.Err
}

func makeScheduleImageBuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return ScheduleImageBuildResponse{}, nil
	}
}

type GetScheduledImageBuildWorkloadsRequest struct {
}

type GetScheduledImageBuildWorkloadsResponse struct {
	Err error `json:"-"`
}

func (r GetScheduledImageBuildWorkloadsResponse) Failed() error {
	return r.Err
}

func makeGetScheduledImageBuildWorkloadsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return GetScheduledImageBuildWorkloadsResponse{}, nil
	}
}

type GetImageBuildStatusRequest struct {
}

type GetImageBuildStatusResponse struct {
	Err error `json:"-"`
}

func (r GetImageBuildStatusResponse) Failed() error {
	return r.Err
}

func makeGetImageBuildStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return GetImageBuildStatusResponse{}, nil
	}
}
