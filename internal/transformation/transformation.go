package transformation

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/utils"
)

type Transformer = func(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error)

type Transformation struct {
	Name        string
	Broker      pubsub.Broker
	Subscriber  pubsub.Subscriber
	inputs      chan *optimusv1.LogEvent
	transformer Transformer
}

func (t *Transformation) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down transformation event loop")
			return nil
		case event := <-t.inputs:
			if t.transformer != nil {
				var newEvent *optimusv1.LogEvent
				newData, err := t.transformer(ctx, event.Data)
				if err != nil {
					slog.Error("could not transform event", "error", err)
					continue
				} else {
					newEvent, err = utils.CopyLogEvent(event)
					if err != nil {
						slog.Error("could not copy log event", "error", err)
						continue
					}
					newEvent.Data = newData
				}
				if newEvent != nil {
					t.Broker.Broadcast(newEvent)
				}
			}

		}
	}
}

func (t *Transformation) Init() pubsub.Broker {
	t.Broker = pubsub.NewBroker(t.Name)
	t.inputs = make(chan *optimusv1.LogEvent, 1)
	t.Subscriber = pubsub.NewSubscriber(t.Name, t.inputs)
	return t.Broker
}

func New(name string, transformer Transformer) *Transformation {
	return &Transformation{
		Name:        name,
		transformer: transformer,
	}
}
