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

func (l *loggingMiddlware) ScheduleWorkload(ctx context.Context) error {
	panic("implement me")
}
