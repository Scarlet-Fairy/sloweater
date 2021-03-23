package redis

import "github.com/Scarlet-Fairy/sloweater/pkg/service"

type redisRepository struct {

}

func New() service.Repository {
	return redisRepository{}
}
