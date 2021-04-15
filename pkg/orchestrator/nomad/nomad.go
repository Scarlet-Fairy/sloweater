package nomad

import (
	"context"
	"fmt"
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

func (n nomadOrchestrator) ScheduleBuildImageJob(ctx context.Context, jobId service.WorkloadId, githubRepo string) (*string, *string, error) {
	_, err := n.ScheduleBatchJob(
		ctx,
		jobId,
		jobId.NameImageBuild(),
		CoboldImage,
		n.imageBuilderArgs(string(jobId), githubRepo),
		nil,
		[]*api.Service{
			n.imageBuilderServiceSettings(string(jobId)),
		},
	)

	if err != nil {
		return nil, nil, err
	}

	jobName := jobId.NameImageBuild()
	imageName := jobId.ImageName(n.registryEndpoint)

	return &jobName, &imageName, nil
}

func (n nomadOrchestrator) imageBuilderArgs(jobId string, githubRepo string) []string {
	return []string{
		"--git-repo", githubRepo,
		"--job-id", jobId,
		"--dev", "false",
		"--redis-url", fmt.Sprintf("localhost:%s", n.config.Orchestrate.ImageBuilder.Services.RedisServicePort),
		"--registry-url", fmt.Sprintf("http://localhost:%s", n.config.Orchestrate.ImageBuilder.Services.RegistryServicePort),
	}
}

func (n nomadOrchestrator) imageBuilderServiceSettings(id string) *api.Service {
	return &api.Service{
		Name: id,
		Connect: &api.ConsulConnect{
			SidecarService: &api.ConsulSidecarService{
				Proxy: &api.ConsulProxy{
					Upstreams: []*api.ConsulUpstream{
						{
							DestinationName: n.config.Orchestrate.ImageBuilder.Services.RedisServiceName,
							LocalBindPort:   n.config.Orchestrate.ImageBuilder.Services.RedisServicePort,
						},
						{
							DestinationName: n.config.Orchestrate.ImageBuilder.Services.RegistryServiceName,
							LocalBindPort:   n.config.Orchestrate.ImageBuilder.Services.RegistryServicePort,
						},
					},
				},
			},
		},
	}
}

func (n nomadOrchestrator) ScheduleBatchJob(
	_ context.Context,
	jobId service.WorkloadId,
	jobName, imageName string,
	args []string,
	envs map[string]string,
	services []*api.Service,
) (*api.JobRegisterResponse, error) {
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
	taskGroup.Services = services
	taskGroup.Networks = []*api.NetworkResource{
		{
			Mode: "bridge",
		},
	}

	job := api.NewBatchJob(
		string(jobId),
		jobName,
		n.config.Orchestrate.Region,
		n.config.Orchestrate.PriorityBatchJob,
	)
	job.AddTaskGroup(taskGroup)
	job.Datacenters = n.config.Orchestrate.Datacenters

	res, _, err := n.client.Jobs().Register(job, &api.WriteOptions{})

	return res, err
}
