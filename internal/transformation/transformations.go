package transformation

import (
	"context"
	"log/slog"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
)

type Transformer interface {
	Transform(ctx context.Context, event *optimusv1.LogEvent) *optimusv1.LogEvent
}
type Transformation struct {
	Name        string
	broker      *pubsub.Broker
	transformer Transformer
	inputs      chan *optimusv1.LogEvent
}

func (t *Transformation) Process(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down transformation event loop")
			return nil
		case event := <-t.inputs:
			newEvent := t.transformer.Transform(ctx, event)
			if newEvent != nil {
				t.broker.Broadcast(newEvent)
			}
		}
	}
}
