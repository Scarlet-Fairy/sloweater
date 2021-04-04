package service

import "context"

type Repository interface {
	CreateJob(ctx context.Context) (jobId string, err error)
	GetJob(ctx context.Context, jobId string) (job *Job, err error)
	ListJobs(ctx context.Context) (jobIds []string, err error)
	SetJobStep(ctx context.Context, jobId string, step Step, error *string) (err error)
	SetJobStatus(ctx context.Context, jobId string, status Status) (err error)
}
