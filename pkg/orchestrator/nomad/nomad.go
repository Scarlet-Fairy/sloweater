package nomad

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/hashicorp/nomad/api"
)

type nomadOrchestrator struct {
	client           *api.Client
	config           service.Config
	registryEndpoint string
}

const (
	DriverDocker = "docker"
)

func New(client *api.Client, config service.Config, registryEndpoint string) service.Orchestrator {
	return nomadOrchestrator{
		client:           client,
		config:           config,
		registryEndpoint: registryEndpoint,
	}
}

func (n nomadOrchestrator) ScheduleBuildImageJob(ctx context.Context, jobId service.JobId, githubRepo string) error {
	_, err := n.ScheduleBatchJob(
		ctx,
		jobId,
		jobId.NameImageBuild(),
		jobId.ImageName(n.registryEndpoint),
		n.imageBuilderArgs(string(jobId), githubRepo),
		nil,
	)
	if err != nil {
		return service.ErrBuildImageSchedulation
	}

	return nil
}

func (n nomadOrchestrator) imageBuilderArgs(jobId string, githubRepo string) []string {
	return []string{
		"--github-repo", githubRepo,
		"--job-id", jobId,
		"--dev", "false",
	}
}

func (n nomadOrchestrator) ScheduleBatchJob(_ context.Context, jobId service.JobId, jobName, imageName string, args []string, envs map[string]string) (*api.JobRegisterResponse, error) {
	task := api.NewTask(jobName, DriverDocker)
	task.Env = envs
	task.Config = map[string]interface{}{
		"image":      imageName,
		"args":       args,
		"force_pull": true,
		"logging":    loggingConfig(n.config.Orchestrate.Logging.LokiUrl),
	}

	taskGroup := api.NewTaskGroup(jobName, 1)
	taskGroup.AddTask(task)

	job := api.NewBatchJob(string(jobId), jobName, n.config.Orchestrate.Region, n.config.Orchestrate.PriorityBatchJob)
	job.AddTaskGroup(taskGroup)

	res, _, err := n.client.Jobs().Register(job, &api.WriteOptions{})

	return res, err
}
