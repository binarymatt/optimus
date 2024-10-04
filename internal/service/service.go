package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/gen/optimus/v1/optimusv1connect"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
)

var _ optimusv1connect.OptimusLogServiceHandler = (*Service)(nil)

type Service struct {
	name   string
	Broker pubsub.Broker
	testId string
}

func (s *Service) generateId() string {
	if s.testId != "" {
		return s.testId
	}
	return ulid.Make().String()
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
func (s *Service) logEventFromStruct(raw *structpb.Struct, source, name string) *optimusv1.LogEvent {
	id := s.generateId()
	return &optimusv1.LogEvent{
		Id:     id,
		Data:   raw,
		Source: source,
	}
}
func (s *Service) StoreLogEvent(ctx context.Context, req *connect.Request[optimusv1.StoreLogEventRequest]) (*connect.Response[optimusv1.StoreLogEventResponse], error) {
	slog.Info("storing events")
	if err := validate(req.Msg); err != nil {
		return nil, err
	}
	eventIds := []string{}
	for _, st := range req.Msg.GetEvents() {
		event := s.logEventFromStruct(st, "http", s.name)
		slog.Debug("broadcasting event", "event", event)
		s.Broker.Broadcast(event)
		metrics.IncProcessedRecord("http_input", s.name)
		eventIds = append(eventIds, event.Id)
	}
	return connect.NewResponse(&optimusv1.StoreLogEventResponse{
		EventIds: eventIds,
	}), nil
}

func New(broker pubsub.Broker, name string) *Service {
	return &Service{
		name:   name,
		Broker: broker,
	}
}
