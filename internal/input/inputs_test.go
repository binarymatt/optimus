package input

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/binarymatt/optimus/internal/pubsub"
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
func setupInput(t *testing.T) (*Input, *MockedInitializer, *MockedProcessor) {

	processor := new(MockedProcessor)
	initializer := new(MockedInitializer)
	t.Cleanup(func() {
		processor.AssertExpectations(t)
		initializer.AssertExpectations(t)
	})
	return &Input{
		Processor:  processor.Process,
		Initialize: initializer.Initialize,
	}, initializer, processor
}
func TestProcess(t *testing.T) {
	i, _, processor := setupInput(t)
	ctx := context.Background()
	processor.On("Process", ctx).Return(nil).Once()
	must.NoError(t, i.Process(ctx))
}

func TestProcess_Error(t *testing.T) {
	i, _, processor := setupInput(t)
	ctx := context.Background()
	errOops := errors.New("oops")
	processor.On("Process", ctx).Return(errOops).Once()
	must.ErrorIs(t, errOops, i.Process(ctx))
}

func TestInit(t *testing.T) {
	i, init, _ := setupInput(t)
	init.On("Initialize", "testid", mock.AnythingOfType("*pubsub.broker")).Return(nil).Once()
	b, err := i.Init("testid")
	must.NotNil(t, b)
	must.NoError(t, err)
}

func TestInit_Error(t *testing.T) {
	i, init, _ := setupInput(t)
	errOops := errors.New("oops")
	init.On("Initialize", "testid", mock.AnythingOfType("*pubsub.broker")).Return(errOops).Once()
	b, err := i.Init("testid")
	must.Nil(t, b)
	must.ErrorIs(t, errOops, err)
}
