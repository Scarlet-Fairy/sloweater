package service

import "context"

type Orchestrator interface {
	ScheduleBuildImageJob(ctx context.Context, jobId, githubRepo string) error
}
