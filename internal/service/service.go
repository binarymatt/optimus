package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/reflect/protoreflect"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/gen/optimus/v1/optimusv1connect"
	"github.com/binarymatt/optimus/internal/pubsub"
)

type Service struct {
	Broker *pubsub.Broker
}

func validate(msg protoreflect.ProtoMessage) *connect.Error {
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
	}
	return connect.NewResponse(&optimusv1.StoreLogEventResponse{}), nil
}

func New(broker *pubsub.Broker) optimusv1connect.OptimusLogServiceHandler {
	return &Service{
		Broker: broker,
	}
}
