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

func (l *loggingMiddlware) ScheduleImageBuild(ctx context.Context, workloadId, githubRepo string) (jobName *string, imageName *string, err error) {
	defer func() {
		l.logger.Log(
			"method", "ScheduleImageBuild",
			"workloadId", workloadId,
			"githubRepo", githubRepo,
			"jobName", jobName,
			"imageName", imageName,
			"err", err,
		)
	}()

	return l.next.ScheduleImageBuild(ctx, workloadId, githubRepo)
}

func (l *loggingMiddlware) ScheduleWorkload(ctx context.Context, workloadId string, envs map[string]string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "ScheduleWorkload",
			"workloadId", workloadId,
			"envs", envs,
			"err", err,
		)
	}()

	return l.next.ScheduleWorkload(ctx, workloadId, envs)
}

func (l *loggingMiddlware) UnScheduleJob(ctx context.Context, jobId string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "UnScheduleJob",
			"jobId", jobId,
			"err", err,
		)
	}()

	return l.next.UnScheduleJob(ctx, jobId)
}
