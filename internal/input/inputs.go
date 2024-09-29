package input

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/binarymatt/optimus/internal/pubsub"
)

const (
	KindFile = "file"
	KindHttp = "http"
)

var (
	ErrInvalidInput = errors.New("invalid input config")
)

type InputProcessor interface {
	Initialize(id string, broker pubsub.Broker) error
	Process(context.Context) error
}
type Input struct {
	Kind   string   `hcl:"kind,label"`
	ID     string   `hcl:"id,label"`
	Body   hcl.Body `hcl:",remain"`
	impl   InputProcessor
	Broker pubsub.Broker
}

func New(id, kind string, internal InputProcessor) (*Input, error) {
	in := &Input{
		ID:   id,
		Kind: kind,
		impl: internal,
	}
	_, err := in.Init()
	return in, err
}

func (i *Input) Process(ctx context.Context) (err error) {
	return i.impl.Process(ctx)
}
func (in *Input) Init() (pubsub.Broker, error) {
	// in.ID = id
	in.Broker = pubsub.NewBroker(in.ID)
	if err := in.impl.Initialize(in.ID, in.Broker); err != nil {
		slog.Error("could not setup input", "error", err)
		return nil, err
	}
	return in.Broker, nil
}

func HclImpl(kind string, ctx *hcl.EvalContext, body hcl.Body) (InputProcessor, hcl.Diagnostics) {
	slog.Debug("setting up input implementation")
	var impl InputProcessor
	switch kind {
	case KindFile:
		impl = &FileInput{}
	case KindHttp:
		impl = &HTTPInput{}
	default:
		diags := hcl.Diagnostics{}
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "invalid input",
			Detail:   fmt.Sprintf("%s is not a valid input type", kind),
		})
		return nil, diags

	}
	return impl, gohcl.DecodeBody(body, ctx, impl)

}
