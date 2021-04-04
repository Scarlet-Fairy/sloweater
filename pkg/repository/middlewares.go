package repository

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/log"
)

type Middleware func(repository service.Repository) service.Repository

func LoggingMiddlware(logger log.Logger) Middleware {
	return func(repository service.Repository) service.Repository {
		return &loggingMiddleware{
			next:   repository,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.Repository
	logger log.Logger
}

func (l *loggingMiddleware) CreateJob(ctx context.Context) (jobId string, err error) {
	defer func() {
		l.logger.Log("method", "CreateJob", "jobId", jobId, "err", err)
	}()

	return l.next.CreateJob(ctx)
}

func (l *loggingMiddleware) GetJob(ctx context.Context, jobId string) (job *service.Job, err error) {
	defer func() {
		l.logger.Log("method", "GetJob", "jobId", jobId, "job", job, "err", err)
	}()

	return l.next.GetJob(ctx, jobId)
}

func (l *loggingMiddleware) ListJobs(ctx context.Context) (jobIds []string, err error) {
	defer func() {
		l.logger.Log("method", "ListJobs", "jobIds", jobIds, "err", err)
	}()

	return l.next.ListJobs(ctx)
}

func (l *loggingMiddleware) SetJobStep(ctx context.Context, jobId string, step service.Step, error *string) (err error) {
	defer func() {
		l.logger.Log("method", "SetJobStep", "jobId", jobId, "step", step, "error", *error, "err", err)
	}()
	return l.next.SetJobStep(ctx, jobId, step, error)
}

func (l *loggingMiddleware) SetJobStatus(ctx context.Context, jobId string, status service.Status) (err error) {
	defer func() {
		l.logger.Log("method", "SetJobStatus", "jobId", jobId, "status", status, "err", err)
	}()

	return l.next.SetJobStatus(ctx, jobId, status)
}
