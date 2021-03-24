package redis

import "github.com/Scarlet-Fairy/sloweater/pkg/service"

type ImageBuildMessage struct {
	Topic service.Step `json:"topic"`
	Error string       `json:"error"`
}
