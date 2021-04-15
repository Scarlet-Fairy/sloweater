package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	ScheduleImageBuild(ctx context.Context, workloadId, githubRepo string) (jobName *string, imageName *string, err error)
	ScheduleWorkload(ctx context.Context) error
}

type basicService struct {
	orchestrator Orchestrator
}

func NewService(orchestrator Orchestrator, logger log.Logger) Service {
	var service Service
	{
		service = &basicService{
			orchestrator: orchestrator,
		}
		service = LoggingMiddleware(logger)(service)
	}

	return service
}

func (s basicService) ScheduleImageBuild(ctx context.Context, workloadId string, githubRepo string) (*string, *string, error) {
	jobId, imageName, err := s.orchestrator.ScheduleBuildImageJob(ctx, WorkloadId(workloadId), githubRepo)
	if err != nil {
		return nil, nil, err
	}

	return jobId, imageName, nil
}

func (s basicService) ScheduleWorkload(ctx context.Context) error {
	return nil
}
