package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Middleware func(service Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(service Service) Service {
		return &loggingMiddlware{
			next:   service,
			logger: logger,
		}
	}
}

type loggingMiddlware struct {
	next   Service
	logger log.Logger
}

func (l *loggingMiddlware) ScheduleImageBuild(ctx context.Context, workloadId, githubRepo string) (jobId string, err error) {
	defer func() {
		l.logger.Log(
			"method", "ScheduleImageBuild",
			"workloadId", workloadId,
			"githubRepo", githubRepo,
			"jobId", jobId,
			"err", err,
		)
	}()

	return l.next.ScheduleImageBuild(ctx, workloadId, githubRepo)
}

func (l *loggingMiddlware) GetImageBuildStatus(ctx context.Context, jobId string) (job *Job, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetImageBuildStatus",
			"jobId", jobId,
			"job", job,
			"err", err,
		)
	}()

	return l.next.GetImageBuildStatus(ctx, jobId)
}

func (l *loggingMiddlware) GetSchedulesImageBuildWorkloads(ctx context.Context) (jobIds []string, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetScheduledImageBuildWorkloads",
			"jobIds", jobIds,
			"err", err,
		)
	}()

	return l.next.GetSchedulesImageBuildWorkloads(ctx)
}

func (l *loggingMiddlware) ScheduleWorkload(ctx context.Context) error {
	panic("implement me")
}

func (l *loggingMiddlware) GetWorkloadStatus(ctx context.Context) error {
	panic("implement me")
}
