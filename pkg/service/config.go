package service

type Config struct {
	Orchestrate OrchestrateConfig
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
	LokiUrl string
}

type ImageBuilderConfig struct {
	Services ServicesConfig
}

type ServicesConfig struct {
	RedisServiceName    string
	RedisServicePort    int
	RegistryServiceName string
	RegistryServicePort int
	LokiServiceName     string
	LokiServicePort     int
}
