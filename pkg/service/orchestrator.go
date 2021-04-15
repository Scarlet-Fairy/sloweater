package service

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrBuildImageSchedulation = errors.New("")
)

type Orchestrator interface {
	ScheduleBuildImageJob(ctx context.Context, workloadId WorkloadId, githubRepo string) (jobName *string, imageName *string, err error)
}
