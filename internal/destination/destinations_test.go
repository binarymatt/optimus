package destination

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/mocks"
)

func TestInit(t *testing.T) {
	mocked := mocks.NewMockDestinationProcessor(t)
	mocked.EXPECT().Setup().Return(nil).Once()
	d := &Destination{}
	d.WithProcessor(mocked)
	must.Eq(t, "", d.ID)
	must.Eq(t, 0, d.BufferSize)
	must.Nil(t, d.inputs)
	must.Nil(t, d.Subscriber)

	err := d.Init("testing")
	must.NoError(t, err)
}

func TestInit_Error(t *testing.T) {
	testErr := errors.New("oops")
	mockImpl := mocks.NewMockDestinationProcessor(t)
	d := &Destination{
		impl: mockImpl,
	}
	mockImpl.EXPECT().Setup().Return(testErr)
	err := d.Init("testing")
	must.ErrorIs(t, testErr, err)
}

func TestWithProcessor(t *testing.T) {
	mocked := mocks.NewMockDestinationProcessor(t)
	d := &Destination{}
	must.Nil(t, d.impl)
	d.WithProcessor(mocked)
	must.NotNil(t, d.impl)
}

func TestProcess(t *testing.T) {
	d := &Destination{}
	mocked := mocks.NewMockDestinationProcessor(t)
	mocked.EXPECT().Setup().Return(nil).Once()
	mocked.EXPECT().Close().Return(nil).Once()
	d.WithProcessor(mocked)
	err := d.Init("testing")
	must.NoError(t, err)
	event := &optimusv1.LogEvent{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"test": structpb.NewStringValue("val"),
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	mocked.EXPECT().Deliver(ctx, event).Return(nil).Once()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		d.Process(ctx)
		return nil
	})
	eg.Go(func() error {
		d.inputs <- event
		time.Sleep(5 * time.Millisecond)
		cancel()
		return nil
	})
	err = eg.Wait()
	must.NoError(t, err)
}
