package redis

import (
	"context"
	"encoding/json"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisPubSub struct {
	client     *redis.Client
	repository service.Repository
}

func New(client *redis.Client, repository service.Repository) service.PubSub {
	return &redisPubSub{
		client:     client,
		repository: repository,
	}
}

func (r redisPubSub) ListenImageBuildEvents(ctx context.Context, jobId string) error {
	ps := r.client.Subscribe(ctx, ImageBuildChannel(jobId))

	for rawMsg := range ps.Channel() {
		var msg ImageBuildMessage
		if err := json.Unmarshal([]byte(rawMsg.Payload), &msg); err != nil {
			if err := r.handleErrorInImageBuildJob(ctx, jobId); err != nil {
				return err
			}

			return service.ErrInvalidImageBuildMessage
		}

		if !msg.Topic.IsValid() {
			if err := r.handleErrorInImageBuildJob(ctx, jobId); err != nil {
				return err
			}

			return service.ErrInvalidImageBuildMessage
		}

		if msg.Error != "" {
			if err := r.handleErrorInImageBuildJob(ctx, jobId); err != nil {
				return err
			}

			return errors.New(msg.Error)
		}

		if err := r.repository.SetJobStep(ctx, jobId, msg.Topic); err != nil {
			return err
		}

		if msg.Topic == service.StepPush {
			if err := r.repository.SetJobStatus(ctx, jobId, service.StatusCompleted); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (r redisPubSub) handleErrorInImageBuildJob(ctx context.Context, jobId string) error {
	return r.repository.SetJobStatus(ctx, jobId, service.StatusError)
}
