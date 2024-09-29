package transformation

import (
	"context"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/mocks"
)

func TestNew(t *testing.T) {
	mocked := mocks.NewMockTransformerImpl(t)
	mocked.EXPECT().Initialize().Return(nil)
	tr, err := New("test_name", "test", []string{}, mocked)
	must.NoError(t, err)
	must.NotNil(t, tr.impl)
}

func TestProcess_HappyPath(t *testing.T) {
	mocked := mocks.NewMockTransformerImpl(t)
	ctx, cancel := context.WithCancel(context.Background())
	mocked.EXPECT().Initialize().Return(nil)
	tr, _ := New("test_name", "test", []string{}, mocked)
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return tr.Process(ctx)
	})
	data, err := structpb.NewStruct(map[string]any{})
	mocked.On("Transform", ctx, data).Return(data, nil).Once()
	must.NoError(t, err)
	evt := &optimusv1.LogEvent{
		Id:   "test",
		Data: data,
	}
	eg.Go(func() error {
		tr.Subscriber.Signal(evt)
		time.Sleep(10 * time.Millisecond)
		cancel()
		return nil
	})
	err = eg.Wait()
	must.NoError(t, err)

}
