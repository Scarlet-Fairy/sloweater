package service

import "context"

type Repository interface {
	CreateJob(ctx context.Context) (string, error)
	GetJob(ctx context.Context, jobId string) (*Job, error)
	ListJobs(ctx context.Context) ([]string, error)
	SetJobStep(ctx context.Context, jobId string, step Step, error *string) error
	SetJobStatus(ctx context.Context, jobId string, status Status) error
}
