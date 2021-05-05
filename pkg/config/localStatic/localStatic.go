package localStatic

import "github.com/Scarlet-Fairy/sloweater/pkg/service"

func NewConfig() *service.Config {
	return &service.Config{
		Orchestrate: service.OrchestrateConfig{
			Logging: service.LoggingConfig{
				LokiUrl: "localhost:3100",
			},
			Region:              "",
			PriorityBatchJob:    50,
			PriorityWorkloadJob: 50,
			Datacenters: []string{
				"dc1",
			},
			ImageBuilder: service.ImageBuilderConfig{
				Services: service.ServicesConfig{
					RedisServiceName:    "redis",
					RedisServicePort:    10001,
					RegistryServiceName: "image-registry",
					RegistryServicePort: 5000,
				},
			},
		},
	}
}
