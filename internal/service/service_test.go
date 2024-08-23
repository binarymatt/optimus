package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/testing/protocmp"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/testutil"
	"github.com/binarymatt/optimus/internal/utils"
)

func TestStoreLogEvent(t *testing.T) {
	ev := testutil.BuildTestEvent()
	ch := make(chan *optimusv1.LogEvent)
	broker := pubsub.NewBroker("test")
	broker.AddSubscriber(pubsub.NewSubscriber("destination", ch))
	s := New(broker, "test")
	_, err := s.StoreLogEvent(context.Background(), connect.NewRequest(&optimusv1.StoreLogEventRequest{
		Key:    "testing",
		Events: []*optimusv1.LogEvent{ev},
	}))
	must.NoError(t, err)
	storedEv := <-ch
	expectedEvent, _ := utils.CopyLogEvent(ev)
	expectedEvent.Upstreams = []string{"test"}
	must.Eq(t, expectedEvent, storedEv, must.Cmp(protocmp.Transform()))
}
