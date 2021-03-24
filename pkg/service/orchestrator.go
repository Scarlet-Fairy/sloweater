package service

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrBuildImageSchedulation = errors.New("")
)

type Orchestrator interface {
	ScheduleBuildImageJob(ctx context.Context, jobId JobId, githubRepo string) error
}
