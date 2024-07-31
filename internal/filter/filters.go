package filter

import (
	"context"
	"errors"
	"log/slog"

	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
)

type FilterFunc = func(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)
type Filter struct {
	id            string
	Broker        *pubsub.Broker
	Subscriber    *pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	process       FilterFunc
	Kind          string   `yaml:"kind"`
	Subscriptions []string `yaml:"subscriptions"`
	BufferSize    int      `yaml:"buffer_size"`
}

func (f *Filter) Init(id string) {
	f.id = id
	if f.BufferSize == 0 {
		f.BufferSize = 5
	}
	f.Broker = pubsub.NewBroker(f.id)
	f.inputs = make(chan *optimusv1.LogEvent, f.BufferSize)
	f.Subscriber = pubsub.NewSubscriber(f.id, f.inputs)
}
func (f *Filter) SetupInternal() error {
	return nil
}
func (f *Filter) UnmarshalYAML(n *yaml.Node) error {
	//type I Input
	//slog.Info("inside", "node", n, "input", i)
	for i := 0; i < len(n.Content)/2; i += 2 {
		key := n.Content[i]
		value := n.Content[i+1]
		if key.Kind == yaml.ScalarNode && key.Value == "kind" {
			if value.Kind != yaml.ScalarNode {
				return errors.New("kind is not scalar")
			}
			f.Kind = value.Value
		}
	}
	switch f.Kind {
	case "bexpr":
		var bFilter BexprFilter
		if err := n.Decode(&bFilter); err != nil {
			return err
		}
		if err := bFilter.Setup(); err != nil {
			return err
		}
		f.process = bFilter.Process
	case "quamina":
		var qFilter QuaminaFilter
		if err := n.Decode(&qFilter); err != nil {
			return err
		}
		if err := qFilter.Setup(); err != nil {
			return err
		}
		f.process = qFilter.Process
	}
	return nil
}

func (f *Filter) Process(ctx context.Context) error {
	slog.Debug("starting filter loop", "id", f.id, "type", f.Kind)
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down filter event loop")
			return nil
		case event := <-f.inputs:
			newEvent, err := f.process(ctx, event)
			if err != nil {
				slog.Error("error delivering record - TODO What now with this record", "error", err)
				continue
			}
			if newEvent != nil {
				f.Broker.Broadcast(newEvent)
			}

		}
	}

}
