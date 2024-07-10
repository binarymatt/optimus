package input

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"

	"github.com/binarymatt/optimus/gen/optimus/v1/optimusv1connect"
	"github.com/binarymatt/optimus/internal/logging"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/service"
)

type HTTPInput struct {
	Address string `yaml:"address"`
}

func (hi *HTTPInput) Setup(ctx context.Context, broker *pubsub.Broker) error {
	path, handler := optimusv1connect.NewOptimusLogServiceHandler(service.New(broker),
		connect.WithInterceptors(
			logging.NewInterceptor(),
			service.MetricsInterceptor(),
		),
	)
	slog.Debug("setting up http input", "path", path)
	http.Handle(path, handler)
	return nil
}

func (hi *HTTPInput) Process(ctx context.Context) error {
	return nil
}
