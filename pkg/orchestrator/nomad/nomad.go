package nomad

import "github.com/Scarlet-Fairy/sloweater/pkg/service"

type nomadOrchestrator struct {

}

func New() service.Orchestrator {
	return nomadOrchestrator{}
}
