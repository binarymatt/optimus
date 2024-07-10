package destination

import (
	"context"
	"errors"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var (
	ErrMissingChannel = errors.New("destination cfg is missing channel")
)

type ChannelDestination struct {
	output chan<- *optimusv1.LogEvent
}

func (cd *ChannelDestination) Setup(cfg map[string]any) error {
	return nil
}

func (cd *ChannelDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	if cd.output == nil {
		return ErrMissingChannel
	}
	cd.output <- event
	return nil
}
