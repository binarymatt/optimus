package filter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
)

var (
	ErrInvalidFilter = errors.New("invalid filter")
)

type FilterProcessor interface {
	Process(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)
}
type FilterFunc = func(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)
type Filter struct {
	id            string
	Broker        pubsub.Broker
	Subscriber    pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	process       FilterFunc
	Kind          string   `yaml:"kind"`
	Subscriptions []string `yaml:"subscriptions"`
	BufferSize    int      `yaml:"buffer_size"`
}

func New(id, kind string, subscriptions []string, impl FilterProcessor) *Filter {
	return &Filter{
		id:      id,
		Kind:    kind,
		process: impl.Process,
	}
}

func (f *Filter) Init(id string) pubsub.Broker {
	f.id = id
	if f.BufferSize == 0 {
		f.BufferSize = 5
	}
	f.Broker = pubsub.NewBroker(f.id)
	f.inputs = make(chan *optimusv1.LogEvent, f.BufferSize)
	f.Subscriber = pubsub.NewSubscriber(f.id, f.inputs)
	return f.Broker
}
func (f *Filter) SetupInternal() error {
	return nil
}
func (f *Filter) UnmarshalYAML(n *yaml.Node) error {
	type alias Filter
	tmp := (*alias)(f)
	if err := n.Decode(&tmp); err != nil {
		slog.Error("inside top level decode", "error", err)
		return err
	}
	f.Kind = tmp.Kind
	f.Subscriptions = tmp.Subscriptions
	f.BufferSize = tmp.BufferSize
	switch f.Kind {
	case "bexpr":
		var bFilter BexprFilter
		if err := n.Decode(&bFilter); err != nil {
			fmt.Println(err)
			return err
		}
		if err := bFilter.Setup(); err != nil {
			fmt.Println("error during setup")
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
	default:
		return ErrInvalidFilter
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
			slog.Warn("broadcasting from filter...")
			if newEvent != nil {
				f.Broker.Broadcast(newEvent)
			}

		}
	}

}
