package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/testutil"
	"github.com/binarymatt/optimus/internal/utils"
)

func TestStoreLogEvent(t *testing.T) {
	ev := testutil.BuildTestEvent()
	ev.Source = "http"
	ch := make(chan *optimusv1.LogEvent)
	broker := pubsub.NewBroker("test")
	broker.AddSubscriber(pubsub.NewSubscriber("destination", ch))
	s := New(broker, "test")
	s.testId = "test"
	_, err := s.StoreLogEvent(context.Background(), connect.NewRequest(&optimusv1.StoreLogEventRequest{
		Key:    "testing",
		Events: []*structpb.Struct{ev.Data},
	}))
	must.NoError(t, err)
	storedEv := <-ch
	expectedEvent, _ := utils.CopyLogEvent(ev)
	expectedEvent.Upstreams = []string{"test"}
	must.Eq(t, expectedEvent, storedEv, must.Cmp(protocmp.Transform()))
}
