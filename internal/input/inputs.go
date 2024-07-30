package input

import (
	"context"
	"errors"
	"log/slog"

	"gopkg.in/yaml.v3"

	"github.com/binarymatt/optimus/internal/pubsub"
)

var (
	ErrInvalidInput = errors.New("invalid input config for kind")
)

type Processor = func(context.Context) error
type InternalSetup = func(id string, broker *pubsub.Broker) error

type Input struct {
	ID        string
	Kind      string `yaml:"kind"`
	Setup     InternalSetup
	Processor Processor
	Broker    *pubsub.Broker
}

func (i *Input) Process(ctx context.Context) (err error) {
	return i.Processor(ctx)
}

func (in *Input) Init(id string) {
	in.ID = id
	in.Broker = pubsub.NewBroker(in.ID)
	if err := in.Setup(in.ID, in.Broker); err != nil {
		slog.Error("could not setup input", "error", err)
	}
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
		in.Setup = finput.Setup
		in.Processor = finput.Process
	case "http":
		var hin HTTPInput
		if err := n.Decode(&hin); err != nil {
			return err
		}
		in.Setup = hin.Setup
		in.Processor = hin.Process
	}
	return nil
}
