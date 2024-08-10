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
	Output chan<- *optimusv1.LogEvent
}

func (cd *ChannelDestination) Setup() error {
	return nil
}

func (cd *ChannelDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	if cd.Output == nil {
		return ErrMissingChannel
	}
	cd.Output <- event
	return nil
}
func (cd *ChannelDestination) Close() error {
	close(cd.Output)
	return nil
}
