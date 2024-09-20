package destination

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
)

var (
	ErrNoProcessor = errors.New("no internal processor")
)

type Deliverer = func(context.Context, *optimusv1.LogEvent) error
type Initializer = func() error
type Closer = func() error
type DestinationProcessor interface {
	Setup() error
	Deliver(context.Context, *optimusv1.LogEvent) error
	Close() error
}
type Destination struct {
	ID            string
	Kind          string   `yaml:"kind"`
	BufferSize    int      `yaml:"buffer_size"`
	Subscriptions []string `yaml:"subscriptions"`
	Subscriber    pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	//process       Deliverer
	//Initialize    Initializer
	//closer        Closer
	impl DestinationProcessor
	// InternalConfig map[string]any `yaml:",inline"`
}

func New(id, kind string, subscriptions []string, impl DestinationProcessor) *Destination {
	d := &Destination{
		ID:            id,
		Kind:          kind,
		Subscriptions: subscriptions,
	}
	d.WithProcessor(impl)
	return d
}

func (d *Destination) WithProcessor(impl DestinationProcessor) {
	slog.Info("setting up processor", "id", d.ID, "kind", d.Kind)
	d.impl = impl
}

func (d *Destination) UnmarshalYAML(n *yaml.Node) error {
	type alias Destination
	tmp := (*alias)(d)
	if err := n.Decode(&tmp); err != nil {
		return err
	}

	d.Kind = tmp.Kind
	d.Subscriptions = tmp.Subscriptions
	d.BufferSize = tmp.BufferSize
	var internal DestinationProcessor
	switch d.Kind {
	case "stdout":
		var std StdOutDestination
		if err := n.Decode(&std); err != nil {
			return err
		}
		internal = &std
	case "http":
		var hout HttpDestination
		if err := n.Decode(&hout); err != nil {
			return err
		}
		internal = &hout
	case "file":
		var fout FileDestination
		if err := n.Decode(&fout); err != nil {
			return err
		}
		internal = &fout
	}
	if internal == nil {
		return ErrNoProcessor
	}
	d.WithProcessor(internal)
	return nil
}

func (d *Destination) Init(id string) error {
	d.ID = id
	slog.Info("initializaing destination", "id", d.ID, "subscriptions", d.Subscriptions, "kind", d.Kind)
	if d.BufferSize == 0 {
		d.BufferSize = 5
	}
	if d.inputs == nil {
		d.inputs = make(chan *optimusv1.LogEvent, d.BufferSize)
	}
	if d.Subscriber == nil {
		d.Subscriber = pubsub.NewSubscriber(d.ID, d.inputs)
	}
	return d.impl.Setup()
}

func (d *Destination) Process(ctx context.Context) {
	defer func() {
		if err := d.impl.Close(); err != nil {
			slog.ErrorContext(ctx, "close destination error", "error", err, "id", d.ID, "kind", d.Kind)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "context is done, shutting down destination event loop", "id", d.ID, "kind", d.Kind)
			return
		case event := <-d.inputs:
			slog.Debug("delivering event", "event", event, "kind", d.Kind, "id", d.ID)
			err := d.impl.Deliver(ctx, event)
			if err != nil {
				slog.Error("error delivering record", "error", err)
			}
			metrics.IncProcessedRecord(fmt.Sprintf("%s_destination", d.Kind), d.ID)
		}
	}

}
