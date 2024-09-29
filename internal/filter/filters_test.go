package filter

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/mocks"
)

func TestInit(t *testing.T) {
	mocked := mocks.NewMockFilterProcessor(t)
	mocked.EXPECT().Setup().Return(nil)
	f := &Filter{
		ID:   "test",
		impl: mocked,
	}
	must.Nil(t, f.Broker)
	must.Nil(t, f.Subscriber)
	must.Nil(t, f.inputs)
	must.Zero(t, f.BufferSize)

	_, err := f.Init()
	must.NoError(t, err)
	must.Eq(t, 5, f.BufferSize)
	must.NotNil(t, f.Broker)
	must.NotNil(t, f.inputs)
	must.NotNil(t, f.Subscriber)
}

func TestUnknownFilter(t *testing.T) {
	t.SkipNow()
	// f := &Filter{}
}

func TestProcess(t *testing.T) {
	event := &optimusv1.LogEvent{
		Id: "test",
	}
	mocked := mocks.NewMockFilterProcessor(t)
	mocked.EXPECT().Setup().Return(nil)
	f := &Filter{
		BufferSize: 1,
		impl:       mocked,
	}
	_, err := f.Init()
	must.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mocked.EXPECT().Process(ctx, event).Return(event, nil).Once()
	eg.Go(func() error {
		// send event to channel
		f.inputs <- event
		// cancel context
		cancel()
		return nil
	})
	eg.Go(func() error {
		return f.Process(ctx)
	})
	must.NoError(t, eg.Wait())

}
