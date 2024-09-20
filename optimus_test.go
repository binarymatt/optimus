package optimus

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/binarymatt/optimus/config"
	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/mocks"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	o, err := New(cfg)
	must.NoError(t, err)
	must.NotNil(t, o)
}

func TestSetup_Inputs(t *testing.T) {
	cfg := config.New()
	cfg.Inputs["test"] = &input.Input{
		Kind: "http",
		Initialize: func(id string, broker pubsub.Broker) error {
			return nil
		},
	}
	o := &Optimus{
		cfg:     cfg,
		parents: make(map[string]pubsub.Broker),
	}
	must.NoError(t, o.setup())
	must.True(t, o.cfg.HttpInputEnabled)
	b, ok := o.parents["test"]
	must.True(t, ok)
	must.NotNil(t, b)
}

func TestSetup_InputError(t *testing.T) {
	cfg := config.New()
	errOops := errors.New("oops")
	cfg.Inputs["test"] = &input.Input{
		Kind: "http",
		Initialize: func(id string, broker pubsub.Broker) error {
			return errOops
		},
	}
	o := &Optimus{
		cfg:     cfg,
		parents: make(map[string]pubsub.Broker),
	}
	must.ErrorIs(t, o.setup(), errOops)
}

func TestSetup_Filters(t *testing.T) {
	cfg := config.New()
	cfg.Filters["test"] = &filter.Filter{
		Kind: "test_filter",
	}
	o := &Optimus{
		cfg:     cfg,
		parents: make(map[string]pubsub.Broker),
	}
	must.NoError(t, o.setup())
	b, ok := o.parents["test"]
	must.True(t, ok)
	must.NotNil(t, b)
}
func TestSetup_Destinations(t *testing.T) {
	cfg := config.New()
	mockDestImpl := mocks.NewMockDestinationProcessor(t)
	mockDestImpl.EXPECT().Setup().Return(nil)
	d := destination.New("test", "http", []string{}, mockDestImpl)
	cfg.Destinations["test"] = d
	o := &Optimus{
		cfg:     cfg,
		parents: make(map[string]pubsub.Broker),
	}
	must.NoError(t, o.setup())
}

func TestSetup_DestinationError(t *testing.T) {
	cfg := config.New()
	errOops := errors.New("oops")
	mockDestImpl := mocks.NewMockDestinationProcessor(t)
	mockDestImpl.EXPECT().Setup().Return(errOops)
	d := destination.New("test", "http", []string{}, mockDestImpl)
	cfg.Destinations["test"] = d
	o := &Optimus{
		cfg:     cfg,
		parents: make(map[string]pubsub.Broker),
	}
	must.ErrorIs(t, o.setup(), errOops)
}
