package orchestrator

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/log"
)

type Middleware func(orchestrator service.Orchestrator) service.Orchestrator

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(orchestrator service.Orchestrator) service.Orchestrator {
		return &loggingMiddleware{
			next:   orchestrator,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.Orchestrator
	logger log.Logger
}

func (l *loggingMiddleware) ScheduleBuildImageJob(ctx context.Context, jobId service.JobId, githubRepo string) (err error) {
	defer func() {
		l.logger.Log("method", "ScheduleBuildImageJob", "jobId", jobId, "githubRepo", githubRepo, "err", err)
	}()

	return l.next.ScheduleBuildImageJob(ctx, jobId, githubRepo)
}
