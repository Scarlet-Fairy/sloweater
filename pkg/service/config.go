package service

type Config struct {
	Orchestrate OrchestrateConfig
	Ingress     IngressConfig
}

type OrchestrateConfig struct {
	Logging             LoggingConfig
	Region              string
	PriorityBatchJob    int
	PriorityWorkloadJob int
	Datacenters         []string
	RestartAttemps      int
	ImageBuilder        ImageBuilderConfig
}

type LoggingConfig struct {
	ElasticUrl string
}

type ImageBuilderConfig struct {
	Services ServicesConfig
}

type ServicesConfig struct {
	RabbitMQServiceName string
	RabbitMQServicePort int

	RegistryServiceName string
	RegistryServicePort int

	ElasticServiceName string
	ElasticServicePort int
}

type IngressConfig struct {
	Protocol string
	Name     string
	Host     string
	Port     int
}
