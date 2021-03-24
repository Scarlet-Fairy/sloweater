package nomad

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/hashicorp/nomad/api"
)

type nomadOrchestrator struct {
	client *api.Client
	region string
}

const (
	PriorityImageBuild = 50
	DriverDocker       = "docker"
)

func New(client *api.Client, region string) service.Orchestrator {
	return nomadOrchestrator{
		client: client,
		region: region,
	}
}

func (n nomadOrchestrator) ScheduleBuildImageJob(ctx context.Context, jobId service.JobId, githubRepo string) error {
	_, err := n.ScheduleBatchJob(ctx, jobId, jobId.NameImageBuild(), jobId.ImageName(nil), nil, nil)
	if err != nil {
		return service.ErrBuildImageSchedulation
	}
}

func (n nomadOrchestrator) ScheduleBatchJob(ctx context.Context, jobId service.JobId, jobName, imageName string, args []string, envs map[string]string) (*api.JobRegisterResponse, error) {
	task := api.NewTask(jobName, DriverDocker)
	task.Env = envs

	taskGroup := api.NewTaskGroup(jobName, 1)
	taskGroup.AddTask(task)

	job := api.NewBatchJob(string(jobId), jobName, n.region, PriorityImageBuild)
	job.AddTaskGroup(taskGroup)

	res, _, err := n.client.Jobs().Register(job, &api.WriteOptions{})

	return res, err
}
