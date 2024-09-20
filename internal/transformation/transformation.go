package transformation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/utils"
)

var (
	ErrUnknownTransformer = errors.New("unkown transformer")
)

type TransformerImpl interface {
	Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error)
	Initialize() error
}
type Transformer = func(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error)
type Initializer = func() error
type Transformation struct {
	ID            string
	Kind          string   `yaml:"kind"`
	Subscriptions []string `yaml:"subscriptions"`
	Broker        pubsub.Broker
	Subscriber    pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	// transformer   Transformer
	// initialize    Initializer
	impl TransformerImpl
}

func (t *Transformation) UnmarshalYAML(n *yaml.Node) error {
	type alias Transformation
	tmp := (*alias)(t)
	if err := n.Decode(&tmp); err != nil {
		return err
	}
	t.Kind = tmp.Kind
	t.Subscriptions = tmp.Subscriptions
	var impl TransformerImpl
	switch t.Kind {
	case "jmespath":
		var jmes JmesTransformer
		if err := n.Decode(&jmes); err != nil {
			return err
		}
		impl = &jmes
	case "jsonata":
		var jsonata JsonataTransformer
		if err := n.Decode(&jsonata); err != nil {
			return err
		}
		impl = &jsonata
	default:
		return fmt.Errorf("%s transformer bad: %w", t.Kind, ErrUnknownTransformer)
	}
	// t.transformer = impl.Transform
	t.impl = impl
	return nil
}
func (t *Transformation) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down transformation event loop")
			return nil
		case event := <-t.inputs:
			var newEvent *optimusv1.LogEvent
			newData, err := t.impl.Transform(ctx, event.Data)
			if err != nil {
				slog.Error("could not transform event", "error", err)
				continue
			} else {
				if newData == nil {
					return nil
				}
				newEvent, err = utils.CopyLogEvent(event)
				if err != nil {
					slog.Error("could not copy log event", "error", err)
					continue
				}
				newEvent.Data = newData
			}
			if newEvent != nil {
				slog.Info("broadcasting")
				t.Broker.Broadcast(newEvent)
			}

		}
	}
}

func (t *Transformation) Init(id string) (pubsub.Broker, error) {
	t.ID = id
	slog.Warn("initializing transformer", "id", t.ID)
	t.Broker = pubsub.NewBroker(t.ID)
	t.inputs = make(chan *optimusv1.LogEvent, 1)
	t.Subscriber = pubsub.NewSubscriber(t.ID, t.inputs)
	return t.Broker, t.impl.Initialize()
}

func New(id, kind string, subscriptions []string, transformer TransformerImpl) *Transformation {
	return &Transformation{
		ID:   id,
		Kind: kind,
		impl: transformer,
	}
}
