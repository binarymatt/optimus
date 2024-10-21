package input

import (
	"context"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/mocks"
)

func TestChannelInitialize(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	ci := &ChannelInput{
		Input: ch,
	}
	broker := pubsub.NewBroker("testname")
	err := ci.Initialize("testname", broker)
	must.NoError(t, err)
	must.Eq(t, "testname", ci.id)
	must.NotNil(t, ci.broker)
}
func TestChannelProcess(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	broker := mocks.NewMockBroker(t)
	evt := &optimusv1.LogEvent{}
	ci := &ChannelInput{
		id:     "test",
		broker: broker,
		Input:  ch,
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	broker.EXPECT().Broadcast(evt).Once()
	eg.Go(func() error {
		return ci.Process(ctx)
	})
	eg.Go(func() error {
		ch <- evt
		time.Sleep(10 * time.Millisecond)
		cancel()
		return nil
	})
	must.NoError(t, eg.Wait())

}
