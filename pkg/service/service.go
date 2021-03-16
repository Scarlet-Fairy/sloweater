package service

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	ScheduleImageBuild() (error)
	ScheduleWorkload()
	GetImageBuildStatus()
	GetWorkloadStatus()
}

func NewService() Service {
	return &sloweaterService{}
}

type sloweaterService struct {

}

