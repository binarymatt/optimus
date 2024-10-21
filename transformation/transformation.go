package transformation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/utils"
)

const (
	KindJmespath = "jmespath"
	KindJsonata  = "jsonata"
)

var (
	ErrUnknownTransformer = errors.New("unkown transformer")
)

type TransformerImpl interface {
	Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error)
	Initialize() error
}
type Transformation struct {
	ID            string
	Kind          string
	Subscriptions []string
	Broker        pubsub.Broker
	Subscriber    pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	impl          TransformerImpl
	BufferSize    int
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

func (t *Transformation) Init() (pubsub.Broker, error) {
	slog.Warn("initializing transformer", "id", t.ID)
	if t.BufferSize == 0 {
		t.BufferSize = 5
	}

	t.Broker = pubsub.NewBroker(t.ID)
	t.inputs = make(chan *optimusv1.LogEvent, t.BufferSize)
	t.Subscriber = pubsub.NewSubscriber(t.ID, t.inputs)
	return t.Broker, t.impl.Initialize()
}

func New(id, kind string, bufferSize int, subscriptions []string, transformer TransformerImpl) (*Transformation, error) {
	t := &Transformation{
		ID:         id,
		Kind:       kind,
		impl:       transformer,
		BufferSize: bufferSize,
	}
	_, err := t.Init()
	return t, err
}
func HclImpl(kind string, body hcl.Body) (TransformerImpl, hcl.Diagnostics) {
	var impl TransformerImpl
	switch kind {
	case KindJmespath:
		impl = &JmesTransformer{}
	case KindJsonata:
		impl = &JsonataTransformer{}
	default:
		diags := append(hcl.Diagnostics{}, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "invalid transformation",
			Detail:   fmt.Sprintf("%s is not a valid transformation type", kind),
		})
		return nil, diags
	}
	return impl, gohcl.DecodeBody(body, nil, impl)

}
