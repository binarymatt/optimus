package config

import (
	"os"
	"testing"

	"github.com/shoenig/test/must"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
	"github.com/binarymatt/optimus/internal/transformation"
	"github.com/binarymatt/optimus/mocks"
)

var ymlStr = `---
data_dir: "/tmp/data"
metrics_enabled: true
listen_address: ":8080"
console: true
log_level: debug
inputs:
  fileInput:
    kind: file
    path: "./cmd/test/tmp"
  httpInput:
    kind: http
destinations:
  sampleout:
    kind: stdout
    subscriptions:
      - fileInput
      - httpInput
      - testing
  samplefile:
    kind: file
    path: "test.ndjson"
    subscriptions:
      - httpInput
`

func TestConfigInit(t *testing.T) {
	cfg := &Config{}
	cfg.Init()
	must.Eq(t, ":8080", cfg.ListenAddress)
}

func TestLoadYamlFile(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	must.NoError(t, err)
	defer os.Remove(f.Name())
	_, err = f.Write([]byte(ymlStr))
	must.NoError(t, err)

	f.Close()
	cfg, err := LoadYamlFromFile(f.Name())
	must.NoError(t, err)

	must.Eq(t, "stdout", cfg.Destinations["sampleout"].Kind)
	must.Eq(t, "file", cfg.Destinations["samplefile"].Kind)
}

func TestLoadYamlFile_NotPresent(t *testing.T) {
	cfg, err := LoadYamlFromFile("oops.path")
	must.Eq(t, "open oops.path: no such file or directory", err.Error())
	must.Nil(t, cfg)
}
func TestLoadYaml(t *testing.T) {
	cfg, err := LoadYaml([]byte(ymlStr))
	must.NoError(t, err)
	must.Eq(t, "stdout", cfg.Destinations["sampleout"].Kind)
	must.Eq(t, "file", cfg.Destinations["samplefile"].Kind)
	//must.Eq(t, cfg.Destinations, expected)
}

func TestWithChannelInput(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	opt := WithChannelInput("testing", ch)
	c := New(opt)
	input, ok := c.Inputs["testing"]
	must.True(t, ok)
	must.NotNil(t, input)
	must.Eq(t, "testing", input.ID)
	must.Eq(t, "channel", input.Kind)
}

func TestWithChannelOutput(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	opt := WithChannelOutput("testOut", ch, []string{"test"})
	cfg := New(opt)
	out, ok := cfg.Destinations["testOut"]
	must.True(t, ok)
	must.NotNil(t, out)
	must.Eq(t, "channel", out.Kind)
	must.Eq(t, []string{"test"}, out.Subscriptions)
}

func TestWithTransformer(t *testing.T) {
	trImpl := mocks.NewMockTransformerImpl(t)
	opt := WithTransformer("testname", "test", []string{}, trImpl)
	cfg := New(opt)
	out, ok := cfg.Transformations["testname"]
	must.True(t, ok)
	must.NotNil(t, out)
}

func TestProgramaticConfig(t *testing.T) {
	cfg := New(
		WithInput("test_input", "http", &input.HTTPInput{}),
		WithFilter("test_filter", "bexpr", []string{"test_input"}, &filter.BexprFilter{Expression: `action == "create"`}),
		WithTransformer("test_transform", "jsonata", []string{"test_filter"}, &transformation.JsonataTransformer{Expression: `{"user_email":principal.email,"path":path}`}),
		WithDestination("test_dest", "stdout", []string{"test_transform"}, &destination.StdOutDestination{}),
		WithMetricsEnabled(),
		WithListenAddress(":8081"),
	)
	must.MapLen(t, 1, cfg.Inputs)
	must.MapLen(t, 1, cfg.Filters)
	must.MapLen(t, 1, cfg.Transformations)
	must.MapLen(t, 1, cfg.Destinations)
	must.Eq(t, ":8081", cfg.ListenAddress)
	must.True(t, cfg.MetricsEnabled)
}

var testYaml = `---
metrics_enabled: true
listen_address: :8081
inputs:
  test_input:
    kind: http
filters:
  test_filter:
    kind: bexpr
    expression: action == "create"
    subscriptions:
      - test_input
transformations:
  test_transform:
    kind: jsonata
    expression: '{"user_email":principal.email,"path":path}'
    subscriptions:
      - test_filter
destinations:
  test_dest:
    kind: stdout
    subscriptions:
      - testtransform
`

func TestCompareConfig(t *testing.T) {
	cfg := New(
		WithInput("test_input", "http", &input.HTTPInput{}),
		WithFilter("test_filter", "bexpr", []string{"test_input"}, &filter.BexprFilter{Expression: `action == "create"`}),
		WithTransformer("test_transform", "jsonata", []string{"test_filter"}, &transformation.JsonataTransformer{Expression: `{"user_email":principal.email,"path":path}`}),
		WithDestination("test_dest", "stdout", []string{"test_transform"}, &destination.StdOutDestination{}),
		WithMetricsEnabled(),
		WithListenAddress(":8081"),
	)
	yamlCfg, err := NewWithYaml([]byte(testYaml))
	must.NoError(t, err)
	must.Eq(t, yamlCfg.HttpInputEnabled, cfg.HttpInputEnabled)
}
