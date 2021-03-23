package service

import "github.com/hashicorp/go.net/context"

type Orchestrator interface {
	ScheduleBuildImageJob(ctx context.Context, jobId, githubRepo string) error
}
