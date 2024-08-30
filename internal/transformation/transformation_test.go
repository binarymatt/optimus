package transformation

import (
	"context"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type MockedTransformer struct {
	mock.Mock
}

func (m *MockedTransformer) Transform(ctx context.Context, data *structpb.Struct) (*structpb.Struct, error) {
	args := m.Called(ctx, data)
	newData := args.Get(0).(*structpb.Struct)
	return newData, args.Error(1)
}
func setupTransformer(t *testing.T) *MockedTransformer {

	transformer := new(MockedTransformer)
	t.Cleanup(func() {
		transformer.AssertExpectations(t)
	})
	return transformer
}
func TestNew(t *testing.T) {
	mocked := setupTransformer(t)
	tr := New("test_name", nil)
	must.Eq(t, "test_name", tr.Name)
	must.Nil(t, tr.transformer)
	tr = New("test_name", mocked.Transform)
	must.NotNil(t, tr.transformer)
}

func TestProcess_HappyPath(t *testing.T) {
	mocked := setupTransformer(t)
	ctx, cancel := context.WithCancel(context.Background())
	tr := New("test_name", mocked.Transform)
	_ = tr.Init()
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
