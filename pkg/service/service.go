package service

import (
	"context"
)

type Service interface {
	ScheduleImageBuild(ctx context.Context, workloadId, githubRepo string) (string, error)
	GetImageBuildStatus(ctx context.Context, jobId string) error
	GetSchedulesImageBuildWorkloads(ctx context.Context) ([]string, error)
	ScheduleWorkload(context.Context) error
	GetWorkloadStatus(context.Context) error
}

func NewService(orchestrator Orchestrator, repository Repository) Service {
	return &basicService{
		orchestrator: orchestrator,
		repository:   repository,
	}
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

	if err := s.orchestrator.ScheduleBuildImageJob(ctx, jobId, githubRepo); err != nil {
		return "", err
	}

	if err := s.pubsub.ListenImageBuildEvents(jobId); err != nil {
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
