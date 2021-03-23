package service

type PubSub interface {
	ListenImageBuildEvents(jobId string) error
}
