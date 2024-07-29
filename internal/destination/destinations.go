package destination

import (
	"context"
	"fmt"
	"log/slog"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
)

type Deliverer = func(context.Context, *optimusv1.LogEvent) error
type InternalDestination interface {
	Setup(map[string]any) error
	Deliver(context.Context, *optimusv1.LogEvent) error
}
type Destination struct {
	ID            string `yaml:"id"`
	Kind          string `yaml:"kind"`
	BufferSize    int    `yaml:"buffer_size"`
	Cfg           interface{}
	Subscriptions []string `yaml:"subscriptions"`
	Subscriber    *pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	process       Deliverer
	Internal      map[string]any `yaml:",inline"`
}

func (d *Destination) SetupInternal() error {
	var internal InternalDestination
	switch d.Kind {
	case "stdout":
		internal = &StdOutDestination{}
	case "http":
	case "file":
		internal = &FileDestination{}
	}
	if err := internal.Setup(d.Internal); err != nil {
		return err
	}
	d.process = internal.Deliver
	return nil
}
func (d *Destination) Init() error {
	slog.Debug("initializaing destination", "id", d.ID, "subscriptions", d.Subscriptions, "internal", d.Internal, "kind", d.Kind)
	if err := d.SetupInternal(); err != nil {
		return err
	}
	if d.BufferSize == 0 {
		d.BufferSize = 5
	}
	if d.inputs == nil {
		d.inputs = make(chan *optimusv1.LogEvent, d.BufferSize)
	}
	if d.Subscriber == nil {
		d.Subscriber = pubsub.NewSubscriber(d.ID, d.inputs)
	}
	return nil
}

func (d *Destination) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down destination event loop", "id", d.ID, "kind", d.Kind)
			return
		case event := <-d.inputs:
			slog.Debug("delivering event", "event", event, "deliverer", d.process)
			err := d.process(ctx, event)
			if err != nil {
				slog.Error("error delivering record", "error", err)
			}
			metrics.RecordProcessedRecord(fmt.Sprintf("%s_destination", d.Kind), d.ID)
		}
	}

}
