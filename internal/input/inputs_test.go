package input

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/mocks"
)

type MockedProcessor struct {
	mock.Mock
}

func (m *MockedProcessor) Process(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockedInitializer struct {
	mock.Mock
}

func (m *MockedInitializer) Initialize(id string, broker pubsub.Broker) error {
	args := m.Called(id, broker)
	return args.Error(0)
}

func TestProcess(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		impl: mocked,
	}
	ctx := context.Background()
	mocked.EXPECT().Process(ctx).Return(nil).Once()
	must.NoError(t, i.Process(ctx))
}

func TestProcess_Error(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		impl: mocked,
	}
	ctx := context.Background()
	errOops := errors.New("oops")
	mocked.EXPECT().Process(ctx).Return(errOops).Once()
	must.ErrorIs(t, errOops, i.Process(ctx))
}

func TestInit(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		ID:   "testid",
		impl: mocked,
	}
	mocked.EXPECT().
		Initialize("testid", mock.AnythingOfType("*pubsub.broker")).
		Return(nil).Once()
	b, err := i.Init()
	must.NotNil(t, b)
	must.NoError(t, err)
}

func TestInit_Error(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		ID:   "testid",
		impl: mocked,
	}
	errOops := errors.New("oops")
	mocked.EXPECT().
		Initialize("testid", mock.AnythingOfType("*pubsub.broker")).
		Return(errOops).Once()
	b, err := i.Init()
	must.Nil(t, b)
	must.ErrorIs(t, errOops, err)
}
