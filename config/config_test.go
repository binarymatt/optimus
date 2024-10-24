package config

import (
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"

	"github.com/binarymatt/optimus/destination"
	"github.com/binarymatt/optimus/filter"
	"github.com/binarymatt/optimus/input"
	"github.com/binarymatt/optimus/mocks"
	"github.com/binarymatt/optimus/transformation"
)

func TestConfigInit(t *testing.T) {
	cfg := &Config{}
	cfg.Init()
	must.Eq(t, ":8080", cfg.ListenAddress)
}

func TestProgramaticConfig(t *testing.T) {
	mockInput := mocks.NewMockInputProcessor(t)
	mockInput.EXPECT().Initialize("test_input", mock.AnythingOfType("*pubsub.broker")).Return(nil).Once()
	mockFilter := mocks.NewMockFilterProcessor(t)
	mockFilter.EXPECT().Setup().Return(nil).Once()
	mockTransformer := mocks.NewMockTransformerImpl(t)
	mockTransformer.EXPECT().Initialize().Return(nil).Once()
	mockDestination := mocks.NewMockDestinationProcessor(t)
	mockDestination.EXPECT().Setup().Return(nil).Once()
	cfg := New(
		WithInput("test_input", "http", mockInput),
		WithFilter("test_filter", "bexpr", 1, []string{"test_input"}, mockFilter),
		WithTransformer("test_transform", "jsonata", 1, []string{"test_filter"}, mockTransformer),
		WithDestination("test_dest", "stdout", 1, []string{"test_transform"}, mockDestination),
		WithMetricsEnabled(),
		WithListenAddress(":8081"),
	)
	must.SliceLen(t, 1, cfg.Inputs)
	must.SliceLen(t, 1, cfg.Filters)
	must.SliceLen(t, 1, cfg.Transformations)
	must.SliceLen(t, 1, cfg.Destinations)
	must.Eq(t, ":8081", cfg.ListenAddress)
	must.True(t, cfg.MetricsEnabled)
}

func TestCompareConfig(t *testing.T) {
	t.SkipNow()
	cfg := New(
		WithInput("test_input", "http", &input.HTTPInput{}),
		WithFilter("test_filter", "bexpr", 1, []string{"test_input"}, &filter.BexprFilter{Expression: `action == "create"`}),
		WithTransformer("test_transform", "jsonata", 1, []string{"test_filter"}, &transformation.JsonataTransformer{Expression: `{"user_email":principal.email,"path":path}`}),
		WithDestination("test_dest", "stdout", 1, []string{"test_transform"}, &destination.StdOutDestination{}),
		WithMetricsEnabled(),
		WithListenAddress(":8081"),
	)
	fmt.Println(cfg)
}

func TestHclConfig(t *testing.T) {
	const exampleConfig = `
	metrics_enabled = true
	listen_address = ":8082"
	input "http" "test_input" {
	}
	filter "bexpr" "test_filter" {
		expression = "test == 1"
		subscriptions = ["test_input"]
	}
	transformation "jsonata" "testnata" {
		expression = <<EOT
		{
  		\"name\": FirstName,
  		\"mobile\": Phone[type = \"mobile\"].number
		}
		EOT
		subscriptions = ["test_filter"]
	}
	destination "stdout" "test_out" {
		subscriptions = ["test_filter"]
	}
	`
	cfg, err := LoadHCL("config.hcl", []byte(exampleConfig))
	must.NoError(t, err)
	must.True(t, cfg.MetricsEnabled)
	must.True(t, cfg.HttpInputEnabled)
	must.Eq(t, ":8082", cfg.ListenAddress)
	must.SliceLen(t, 1, cfg.Inputs)
	must.SliceLen(t, 1, cfg.Filters)
	filter := cfg.Filters[0]
	must.Eq(t, []string{"test_input"}, filter.Subscriptions)
	must.Eq(t, "test_filter", filter.ID)
	must.Eq(t, "bexpr", filter.Kind)
}
