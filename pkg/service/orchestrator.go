package service

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrScheduleBatchJob    = errors.New("failed to schedule batch job")
	ErrScheduleWorkloadJob = errors.New("failed to schedule workload job")
	ErrDegisterJob         = errors.New("failed to deregister job")
)

type Orchestrator interface {
	ScheduleBuildImageJob(ctx context.Context, workloadId WorkloadId, githubRepo string) (jobName *string, imageName *string, err error)
	ScheduleWorkloadJob(ctx context.Context, workloadId WorkloadId, envs map[string]string) (err error)
	UnScheduleJob(ctx context.Context, jobId string) (err error)
}
