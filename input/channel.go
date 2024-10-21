package input

import (
	"context"
	"errors"
	"fmt"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
)

var (
	ErrMissingChannel = errors.New("input cfg is missing channel")
)

type ChannelInput struct {
	id     string
	broker pubsub.Broker
	Input  <-chan *optimusv1.LogEvent
}

func (ci *ChannelInput) Initialize(id string, broker pubsub.Broker) error {
	ci.broker = broker
	ci.id = id
	return nil
}

func (ci *ChannelInput) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-ci.Input:
			ci.broker.Broadcast(e)
			metrics.IncProcessedRecord(fmt.Sprintf("%s_input", "channel"), ci.id)
		}
	}
}
