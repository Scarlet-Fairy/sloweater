package nomad

import (
	"context"
	"fmt"
	"github.com/Scarlet-Fairy/sloweater/pkg/orchestrator"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-kit/kit/log"
	consulApi "github.com/hashicorp/consul/api"
	nomadApi "github.com/hashicorp/nomad/api"
	"strconv"
	"strings"
)

type nomadOrchestrator struct {
	nomadClient      *nomadApi.Client
	consulClient     *consulApi.Client
	config           service.Config
	registryEndpoint string
}

const (
	DriverDocker = "docker"
	CoboldImage  = "scarletfairy/cobold"
)

func New(
	nomadClient *nomadApi.Client,
	consulClient *consulApi.Client,
	config service.Config,
	logger log.Logger,
	registryEndpoint string,
) service.Orchestrator {
	var orc service.Orchestrator
	{
		orc = &nomadOrchestrator{
			nomadClient:      nomadClient,
			consulClient:     consulClient,
			config:           config,
			registryEndpoint: registryEndpoint,
		}
		orc = orchestrator.LoggingMiddleware(logger)(orc)
	}

	return orc
}

func (n *nomadOrchestrator) ScheduleBuildImageJob(ctx context.Context, jobId service.WorkloadId, githubRepo string) (*string, *string, error) {
	_, err := n.ScheduleBatchJob(
		ctx,
		jobId,
		jobId.NameImageBuild(),
		CoboldImage,
		n.imageBuilderArgs(string(jobId), githubRepo),
		map[string]string{
			//"DEV": "true",
		},
		[]*nomadApi.Service{
			n.imageBuilderServices(jobId.NameService()),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	jobName := jobId.NameImageBuild()
	imageName := jobId.ImageName(n.registryEndpoint)

	return &jobName, &imageName, nil
}

func (n *nomadOrchestrator) imageBuilderArgs(jobId string, githubRepo string) []string {
	return []string{
		"--git-repo", githubRepo,
		"--job-id", jobId,
		"--docker-registry", fmt.Sprintf("localhost:%d", n.config.Orchestrate.ImageBuilder.Services.RegistryServicePort),
		"--rabbitmq-url", fmt.Sprintf("amqp://guest:guest@localhost:%d/", n.config.Orchestrate.ImageBuilder.Services.RabbitMQServicePort),
	}
}

func (n *nomadOrchestrator) imageBuilderServices(id string) *nomadApi.Service {
	return &nomadApi.Service{
		Name: id,
		Connect: &nomadApi.ConsulConnect{
			SidecarService: &nomadApi.ConsulSidecarService{
				Proxy: &nomadApi.ConsulProxy{
					Upstreams: []*nomadApi.ConsulUpstream{
						{
							DestinationName: n.config.Orchestrate.ImageBuilder.Services.RabbitMQServiceName,
							LocalBindPort:   n.config.Orchestrate.ImageBuilder.Services.RabbitMQServicePort,
						},
						{
							DestinationName: n.config.Orchestrate.ImageBuilder.Services.RegistryServiceName,
							LocalBindPort:   n.config.Orchestrate.ImageBuilder.Services.RegistryServicePort,
						},
						{
							DestinationName: n.config.Orchestrate.ImageBuilder.Services.ElasticServiceName,
							LocalBindPort:   n.config.Orchestrate.ImageBuilder.Services.ElasticServicePort,
						},
					},
				},
			},
		},
	}
}

func (n *nomadOrchestrator) ScheduleBatchJob(
	_ context.Context,
	jobId service.WorkloadId,
	jobName, imageName string,
	args []string,
	envs map[string]string,
	services []*nomadApi.Service,
) (*nomadApi.JobRegisterResponse, error) {
	task := nomadApi.NewTask(jobName, DriverDocker)
	task.Env = envs
	task.Config = map[string]interface{}{
		"image": imageName,
		"args":  args,
		//"force_pull": true,
		"volumes": []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
		"logging": []map[string]interface{}{
			loggingConfig(n.config.Orchestrate.Logging.ElasticUrl, string(jobId)),
		},
	}
	task.RestartPolicy = &nomadApi.RestartPolicy{
		Attempts: &n.config.Orchestrate.RestartAttemps,
	}

	taskGroup := nomadApi.NewTaskGroup(jobName, 1)
	taskGroup.AddTask(task)
	taskGroup.ReschedulePolicy = &nomadApi.ReschedulePolicy{
		Attempts: &n.config.Orchestrate.RestartAttemps,
	}
	taskGroup.Services = services
	taskGroup.Networks = []*nomadApi.NetworkResource{
		{
			Mode: "bridge",
		},
	}

	job := nomadApi.NewBatchJob(
		jobName,
		jobName,
		n.config.Orchestrate.Region,
		n.config.Orchestrate.PriorityBatchJob,
	)
	job.AddTaskGroup(taskGroup)
	job.Datacenters = n.config.Orchestrate.Datacenters

	res, _, err := n.nomadClient.Jobs().Register(job, &nomadApi.WriteOptions{})

	return res, err
}

func (n nomadOrchestrator) ScheduleWorkloadJob(_ context.Context, workloadId service.WorkloadId, envs map[string]string) (*string, *string, error) {
	jobName := workloadId.NameWorkload()
	jobPort := n.workloadPort()
	url := fmt.Sprintf("%s.%s", workloadId.NameService(), n.config.Ingress.Host)

	envs["PORT"] = strconv.Itoa(jobPort)

	task := nomadApi.NewTask(workloadId.NameWorkload(), DriverDocker)
	task.Env = envs
	task.Config = map[string]interface{}{
		"image":      workloadId.ImageName(n.registryEndpoint),
		"force_pull": true,
	}
	task.RestartPolicy = &nomadApi.RestartPolicy{
		Attempts: &n.config.Orchestrate.RestartAttemps,
	}

	taskGroup := nomadApi.NewTaskGroup(workloadId.NameWorkload(), 1)
	taskGroup.AddTask(task)
	taskGroup.Networks = []*nomadApi.NetworkResource{
		{
			Mode: "bridge",
		},
	}
	taskGroup.Services = []*nomadApi.Service{
		{
			Name:      workloadId.NameService(),
			PortLabel: strconv.Itoa(jobPort),
			Connect: &nomadApi.ConsulConnect{
				SidecarService: &nomadApi.ConsulSidecarService{},
			},
		},
	}

	job := nomadApi.NewServiceJob(
		jobName,
		jobName,
		n.config.Orchestrate.Region,
		n.config.Orchestrate.PriorityWorkloadJob,
	)
	job.AddTaskGroup(taskGroup)
	job.Datacenters = n.config.Orchestrate.Datacenters

	if _, _, err := n.nomadClient.Jobs().Register(job, &nomadApi.WriteOptions{}); err != nil {
		return nil, nil, err
	}

	if err := n.updateGatewayIngress(); err != nil {
		return nil, nil, err
	}

	return &jobName, &url, nil
}

func (n nomadOrchestrator) workloadPort() int {
	return 8080
}

func (n nomadOrchestrator) updateGatewayIngress() error {

	allServices, _, err := n.consulClient.Catalog().Services(&consulApi.QueryOptions{})
	if err != nil {
		return err
	}

	var serviceNames []string
	for key := range allServices {
		if strings.HasPrefix(key, fmt.Sprintf("%s-", service.ServiceNameWorkload)) {
			if _, _, err := n.consulClient.ConfigEntries().Set(&consulApi.ServiceConfigEntry{
				Kind:     consulApi.ServiceDefaults,
				Name:     key,
				Protocol: "http",
			}, &consulApi.WriteOptions{}); err != nil {
				return err
			}

			serviceNames = append(serviceNames, key)
		}
	}

	group := nomadApi.NewTaskGroup(n.config.Ingress.Name, 1)
	group.Networks = []*nomadApi.NetworkResource{
		{
			Mode: "bridge",
			ReservedPorts: []nomadApi.Port{
				{
					Label: "inbound",
					Value: n.config.Ingress.Port,
					To:    n.config.Ingress.Port,
				},
			},
		},
	}

	var services []*nomadApi.ConsulIngressService
	for _, name := range serviceNames {
		services = append(services, &nomadApi.ConsulIngressService{
			Name:  name,
			Hosts: []string{fmt.Sprintf("%s.%s", name, n.config.Ingress.Host)},
		})
	}

	group.Services = []*nomadApi.Service{
		{
			Name:      n.config.Ingress.Name,
			PortLabel: strconv.Itoa(n.config.Ingress.Port),
			Connect: &nomadApi.ConsulConnect{
				Gateway: &nomadApi.ConsulGateway{
					Proxy: &nomadApi.ConsulGatewayProxy{},
					Ingress: &nomadApi.ConsulIngressConfigEntry{
						Listeners: []*nomadApi.ConsulIngressListener{
							{
								Port:     n.config.Ingress.Port,
								Protocol: "http",
								Services: services,
							},
						},
					},
				},
			},
		},
	}

	job := nomadApi.NewServiceJob(
		n.config.Ingress.Name,
		n.config.Ingress.Name,
		n.config.Orchestrate.Region,
		n.config.Orchestrate.PriorityWorkloadJob,
	)
	job.Datacenters = n.config.Orchestrate.Datacenters
	job.AddTaskGroup(group)

	if _, _, err = n.nomadClient.Jobs().Register(job, &nomadApi.WriteOptions{}); err != nil {
		return err
	}

	return nil
}

func (n nomadOrchestrator) UnScheduleJob(_ context.Context, jobId string) error {
	_, _, err := n.nomadClient.Jobs().Deregister(jobId, true, &nomadApi.WriteOptions{})
	if err != nil {
		return err
	}

	if err := n.updateGatewayIngress(); err != nil {
		return err
	}

	return nil
}
