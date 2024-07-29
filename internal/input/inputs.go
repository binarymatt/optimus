package input

import (
	"context"
	"errors"

	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
)

var (
	ErrInvalidInput = errors.New("invalid input config for kind")
)

type InputSpecific interface {
	Setup(context.Context, *pubsub.Broker) error
	Process(context.Context) error
	SetID(id string)
}
type Ingester = func(ctx context.Context, event *optimusv1.LogEvent) error

type Input struct {
	ID       string
	Kind     string `yaml:"kind"`
	Internal InputSpecific
	Broker   *pubsub.Broker
	// 	ingester Ingester
}

func (i *Input) Process(ctx context.Context) (err error) {
	err = i.Internal.Process(ctx)
	return
}

func (in *Input) Init(id string) {
	in.ID = id
	in.Internal.SetID(id)
	in.Broker = pubsub.NewBroker(in.ID)
}

func (in *Input) UnmarshalYAML(n *yaml.Node) error {
	//type I Input
	//slog.Info("inside", "node", n, "input", i)
	for i := 0; i < len(n.Content)/2; i += 2 {
		key := n.Content[i]
		value := n.Content[i+1]
		if key.Kind == yaml.ScalarNode && key.Value == "kind" {
			if value.Kind != yaml.ScalarNode {
				return errors.New("kind is not scalar")
			}
			in.Kind = value.Value
		}
	}
	switch in.Kind {
	case "file":
		var finput FileInput
		if err := n.Decode(&finput); err != nil {
			return err
		}
		in.Internal = &finput
	case "http":
		var hin HTTPInput
		if err := n.Decode(&hin); err != nil {
			return err
		}
		in.Internal = &hin
	}
	return nil
}
