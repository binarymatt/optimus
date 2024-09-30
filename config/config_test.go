package config

import (
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
	"github.com/binarymatt/optimus/internal/transformation"
	"github.com/binarymatt/optimus/mocks"
)

func TestConfigInit(t *testing.T) {
	cfg := &Config{}
	cfg.Init()
	must.Eq(t, ":8080", cfg.ListenAddress)
}

func TestWithChannelInput(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	opt := WithChannelInput("testing", ch)
	c := New(opt)
	input := c.Inputs[0]
	must.NotNil(t, input)
	must.Eq(t, "testing", input.ID)
	must.Eq(t, "channel", input.Kind)
}

func TestWithChannelOutput(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	opt := WithChannelOutput("testOut", ch, []string{"test"})
	cfg := New(opt)
	out := cfg.Destinations[0]
	must.NotNil(t, out)
	must.Eq(t, "channel", out.Kind)
	must.Eq(t, []string{"test"}, out.Subscriptions)
}

func TestWithTransformer(t *testing.T) {
	trImpl := mocks.NewMockTransformerImpl(t)
	trImpl.EXPECT().Initialize().Return(nil)
	opt := WithTransformer("testname", "test", []string{}, trImpl)
	cfg := New(opt)
	out := cfg.Transformations[0]
	must.NotNil(t, out)
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
		WithFilter("test_filter", "bexpr", []string{"test_input"}, mockFilter),
		WithTransformer("test_transform", "jsonata", []string{"test_filter"}, mockTransformer),
		WithDestination("test_dest", "stdout", []string{"test_transform"}, mockDestination),
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
		WithFilter("test_filter", "bexpr", []string{"test_input"}, &filter.BexprFilter{Expression: `action == "create"`}),
		WithTransformer("test_transform", "jsonata", []string{"test_filter"}, &transformation.JsonataTransformer{Expression: `{"user_email":principal.email,"path":path}`}),
		WithDestination("test_dest", "stdout", []string{"test_transform"}, &destination.StdOutDestination{}),
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

var rawData = `
metrics_enabled = false
`

func TestLoadHcl(t *testing.T) {
	cfg, err := LoadHCL("test.hcl", []byte(rawData))
	must.NoError(t, err)
	must.NotNil(t, cfg)
}
