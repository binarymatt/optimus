package filter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
)

const (
	KindBexpr   = "bexpr"
	KindQuamina = "quamina"
)

var (
	ErrInvalidFilter = errors.New("invalid filter")
)

type FilterProcessor interface {
	Process(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)
	Setup() error
}
type FilterFunc = func(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)
type Filter struct {
	ID            string
	Broker        pubsub.Broker
	Subscriber    pubsub.Subscriber
	inputs        chan *optimusv1.LogEvent
	impl          FilterProcessor
	Kind          string
	Subscriptions []string
	BufferSize    int
}

func New(id, kind string, subscriptions []string, impl FilterProcessor) (*Filter, error) {
	f := &Filter{
		ID:            id,
		Kind:          kind,
		impl:          impl,
		Subscriptions: subscriptions,
	}
	_, err := f.Init()
	return f, err
}

func (f *Filter) Init() (pubsub.Broker, error) {
	if f.BufferSize == 0 {
		f.BufferSize = 5
	}
	f.Broker = pubsub.NewBroker(f.ID)
	f.inputs = make(chan *optimusv1.LogEvent, f.BufferSize)
	f.Subscriber = pubsub.NewSubscriber(f.ID, f.inputs)

	return f.Broker, f.impl.Setup()
}

func (f *Filter) Process(ctx context.Context) error {
	slog.Debug("starting filter loop", "id", f.ID, "type", f.Kind)
	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, shutting down filter event loop")
			return nil
		case event := <-f.inputs:
			newEvent, err := f.impl.Process(ctx, event)
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
func HclImpl(kind string, ctx *hcl.EvalContext, body hcl.Body) (FilterProcessor, hcl.Diagnostics) {
	var impl FilterProcessor
	switch kind {
	case KindBexpr:
		impl = &BexprFilter{}
	case KindQuamina:
		impl = &QuaminaFilter{}
	default:
		diags := hcl.Diagnostics{}
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "invalid filter",
			Detail:   fmt.Sprintf("%s is not a valid filter type.", kind),
		})
		return nil, diags
	}
	return impl, gohcl.DecodeBody(body, ctx, impl)

}
