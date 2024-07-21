package input

import (
	"context"
	"errors"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
)

var (
	ErrMissingChannel = errors.New("input cfg is missing channel")
)

type ChannelInput struct {
	broker *pubsub.Broker
	Input  <-chan *optimusv1.LogEvent
}

func (ci *ChannelInput) Setup(ctx context.Context, broker *pubsub.Broker) error {
	ci.broker = broker
	return nil
}

func (ci *ChannelInput) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-ci.Input:
			ci.broker.Broadcast(e)
		}
	}
}
