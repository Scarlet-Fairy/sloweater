package service

type Config struct {
	Orchestrate OrchestrateConfig
}

type OrchestrateConfig struct {
	Logging          LoggingConfig
	Region           string
	PriorityBatchJob int
	Datacenters      []string
	RestartAttemps   int
}

type LoggingConfig struct {
	LokiUrl string
}
