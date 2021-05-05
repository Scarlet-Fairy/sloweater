package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	ScheduleImageBuild(ctx context.Context, workloadId, githubRepo string) (jobName *string, imageName *string, err error)
	ScheduleWorkload(ctx context.Context, workloadId string, envs map[string]string) (jobName *string, err error)
	UnScheduleJob(ctx context.Context, jobId string) error
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

func (s basicService) ScheduleWorkload(ctx context.Context, workloadId string, envs map[string]string) (*string, error) {
	jobName, err := s.orchestrator.ScheduleWorkloadJob(ctx, WorkloadId(workloadId), envs)
	if err != nil {
		return nil, err
	}

	return jobName, nil
}

func (s basicService) UnScheduleJob(ctx context.Context, jobId string) error {
	err := s.orchestrator.UnScheduleJob(ctx, jobId)
	if err != nil {
		return err
	}

	return nil
}
