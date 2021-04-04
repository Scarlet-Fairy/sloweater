package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	ScheduleImageBuild(ctx context.Context, workloadId, githubRepo string) (jobId string, err error)
	GetImageBuildStatus(ctx context.Context, jobId string) (job *Job, err error)
	GetSchedulesImageBuildWorkloads(ctx context.Context) (jobIds []string, err error)

	ScheduleWorkload(ctx context.Context) error
	GetWorkloadStatus(ctx context.Context) error
}

func NewService(orchestrator Orchestrator, repository Repository, pubSub PubSub, logger log.Logger) Service {
	var service Service
	{
		service = &basicService{
			orchestrator: orchestrator,
			repository:   repository,
			pubsub:       pubSub,
		}
		service = LoggingMiddleware(logger)(service)
	}

	return service
}

type basicService struct {
	orchestrator Orchestrator
	repository   Repository
	pubsub       PubSub
}

func (s basicService) ScheduleImageBuild(ctx context.Context, workloadId string, githubRepo string) (string, error) {
	jobId, err := s.repository.CreateJob(ctx)
	if err != nil {
		return "", err
	}

	if err := s.orchestrator.ScheduleBuildImageJob(ctx, JobId(jobId), githubRepo); err != nil {
		return "", err
	}

	if err := s.pubsub.ListenImageBuildEvents(ctx, jobId); err != nil {
		return "", err
	}

	return jobId, nil
}

func (s basicService) GetImageBuildStatus(ctx context.Context, jobId string) (*Job, error) {
	job, err := s.repository.GetJob(ctx, jobId)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s basicService) GetSchedulesImageBuildWorkloads(ctx context.Context) ([]string, error) {
	jobs, err := s.repository.ListJobs(ctx)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s basicService) ScheduleWorkload(ctx context.Context) error {
	panic("implement me")
}

func (s basicService) GetWorkloadStatus(ctx context.Context) error {
	panic("implement me")
}
