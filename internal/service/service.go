package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/gen/optimus/v1/optimusv1connect"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
)

type Service struct {
	name   string
	Broker *pubsub.Broker
}

func validate(msg *optimusv1.StoreLogEventRequest) *connect.Error {
	validator, err := protovalidate.New()
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	if err := validator.Validate(msg); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	return nil
}

func (s *Service) StoreLogEvent(ctx context.Context, req *connect.Request[optimusv1.StoreLogEventRequest]) (*connect.Response[optimusv1.StoreLogEventResponse], error) {
	slog.Info("storing events")
	if err := validate(req.Msg); err != nil {
		return nil, err
	}
	for _, event := range req.Msg.GetEvents() {
		slog.Debug("broadcasting event", "event", event)
		s.Broker.Broadcast(event)
		metrics.RecordProcessedRecord("http_input", s.name)
	}
	return connect.NewResponse(&optimusv1.StoreLogEventResponse{}), nil
}

func New(broker *pubsub.Broker, name string) optimusv1connect.OptimusLogServiceHandler {
	return &Service{
		name:   name,
		Broker: broker,
	}
}
