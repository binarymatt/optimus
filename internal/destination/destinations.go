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

type Deliverer = func(context.Context, *optimusv1.LogEvent) error
type Initializer = func() error
type InternalDestination interface {
	Setup(map[string]any) error
	Deliver(context.Context, *optimusv1.LogEvent) error
}
type Destination struct {
	id             string
	Kind           string   `yaml:"kind"`
	BufferSize     int      `yaml:"buffer_size"`
	Subscriptions  []string `yaml:"subscriptions"`
	Subscriber     *pubsub.Subscriber
	inputs         chan *optimusv1.LogEvent
	process        Deliverer
	InternalConfig map[string]any `yaml:",inline"`
}

func (d *Destination) UnmarshalYAML(n *yaml.Node) error {
	for i := 0; i < len(n.Content)/2; i += 2 {
		key := n.Content[i]
		value := n.Content[i+1]
		if key.Kind == yaml.ScalarNode && key.Value == "kind" {
			if value.Kind != yaml.ScalarNode {
				return errors.New("kind is not scalar")
			}
			d.Kind = value.Value
		}
	}
	var internal InternalDestination
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
		return errors.New("did not have an internal processor")
	}
	// in.WithInputProcessor(internal)
	return nil
}

func (d *Destination) SetupInternal() error {
	var internal InternalDestination
	switch d.Kind {
	case "stdout":
		internal = &StdOutDestination{}
	//case "http":
	// 	internal = &HttpDestination{}
	case "file":
		internal = &FileDestination{}
	}
	if err := internal.Setup(d.InternalConfig); err != nil {
		return err
	}
	d.process = internal.Deliver
	return nil
}
func (d *Destination) Init(id string) error {
	d.id = id
	slog.Debug("initializaing destination", "id", d.id, "subscriptions", d.Subscriptions, "internal", d.InternalConfig, "kind", d.Kind)
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
		d.Subscriber = pubsub.NewSubscriber(d.id, d.inputs)
	}
	return nil
}

func (d *Destination) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down destination event loop", "id", d.id, "kind", d.Kind)
			return
		case event := <-d.inputs:
			slog.Debug("delivering event", "event", event, "deliverer", d.process)
			err := d.process(ctx, event)
			if err != nil {
				slog.Error("error delivering record", "error", err)
			}
			metrics.RecordProcessedRecord(fmt.Sprintf("%s_destination", d.Kind), d.id)
		}
	}

}
