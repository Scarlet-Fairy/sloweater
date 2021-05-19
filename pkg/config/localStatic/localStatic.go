package localStatic

import "github.com/Scarlet-Fairy/sloweater/pkg/service"

func NewConfig() *service.Config {
	return &service.Config{
		Orchestrate: service.OrchestrateConfig{
			Logging: service.LoggingConfig{
				ElasticUrl: "http://localhost:9200",
			},
			Region:              "",
			PriorityBatchJob:    50,
			PriorityWorkloadJob: 50,
			Datacenters: []string{
				"dc1",
			},
			ImageBuilder: service.ImageBuilderConfig{
				Services: service.ServicesConfig{
					RabbitMQServiceName: "rabbitmq",
					RabbitMQServicePort: 10001,

					RegistryServiceName: "image-registry",
					RegistryServicePort: 5000,

					ElasticServiceName: "elasticsearch-api",
					ElasticServicePort: 9200,
				},
			},
		},
	}
}
