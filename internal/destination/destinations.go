package destination

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
)

const (
	KindHttp    = "http"
	KindFile    = "file"
	KindStdOut  = "stdout"
	KindChannel = "channel"
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
	Kind          string
	BufferSize    int
	Subscriptions []string
	Subscriber    pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	impl          DestinationProcessor
}

func New(id, kind string, subscriptions []string, impl DestinationProcessor) (*Destination, error) {
	d := &Destination{
		ID:            id,
		Kind:          kind,
		Subscriptions: subscriptions,
		impl:          impl,
	}
	return d, d.Init(id)
}

func (d *Destination) Init(id string) error {
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
func HclImpl(kind string, body hcl.Body) (DestinationProcessor, hcl.Diagnostics) {
	var impl DestinationProcessor
	switch kind {
	case KindFile:
		impl = &FileDestination{}
	case KindHttp:
		impl = &HttpDestination{}
	case KindStdOut:
		impl = &StdOutDestination{}
	default:
		diags := append(hcl.Diagnostics{}, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "invalid destination",
			Detail:   fmt.Sprintf("%s is not a valid destination type", kind),
		})
		return nil, diags
	}
	return impl, gohcl.DecodeBody(body, nil, impl)
}
