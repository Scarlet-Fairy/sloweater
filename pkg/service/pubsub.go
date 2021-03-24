package service

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrInvalidImageBuildMessage = errors.New("unable to parse message from image build job")
)

type PubSub interface {
	ListenImageBuildEvents(ctx context.Context, jobId string) error
}
