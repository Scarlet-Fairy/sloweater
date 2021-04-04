package nomad

import (
	"context"
	"github.com/Scarlet-Fairy/sloweater/pkg/orchestrator"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/nomad/api"
)

type nomadOrchestrator struct {
	client           *api.Client
	config           service.Config
	registryEndpoint string
}

const (
	DriverDocker = "docker"
	CoboldImage  = "scarletfairy/cobold"
)

func New(client *api.Client, config service.Config, logger log.Logger, registryEndpoint string) service.Orchestrator {
	var orc service.Orchestrator
	{
		orc = nomadOrchestrator{
			client:           client,
			config:           config,
			registryEndpoint: registryEndpoint,
		}
		orc = orchestrator.LoggingMiddleware(logger)(orc)
	}

	return orc
}

func (n nomadOrchestrator) ScheduleBuildImageJob(ctx context.Context, jobId service.JobId, githubRepo string) error {
	_, err := n.ScheduleBatchJob(
		ctx,
		jobId,
		jobId.NameImageBuild(),
		CoboldImage,
		n.imageBuilderArgs(string(jobId), githubRepo),
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (n nomadOrchestrator) imageBuilderArgs(jobId string, githubRepo string) []string {
	return []string{
		"--git-repo", githubRepo,
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
		/*"logging": []map[string]interface{}{
			loggingConfig(n.config.Orchestrate.Logging.LokiUrl),
		},*/
	}
	task.RestartPolicy = &api.RestartPolicy{
		Attempts: &n.config.Orchestrate.RestartAttemps,
	}

	taskGroup := api.NewTaskGroup(jobName, 1)
	taskGroup.AddTask(task)
	taskGroup.ReschedulePolicy = &api.ReschedulePolicy{
		Attempts: &n.config.Orchestrate.RestartAttemps,
	}

	job := api.NewBatchJob(string(jobId), jobName, n.config.Orchestrate.Region, n.config.Orchestrate.PriorityBatchJob)
	job.AddTaskGroup(taskGroup)
	job.Datacenters = n.config.Orchestrate.Datacenters

	res, _, err := n.client.Jobs().Register(job, &api.WriteOptions{})

	return res, err
}
