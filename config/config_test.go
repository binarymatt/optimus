package config

import (
	"context"
	"os"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
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
	opt := WithChannelOutput("testOut", ch)
	cfg := New(opt)
	out, ok := cfg.Destinations["testOut"]
	must.True(t, ok)
	must.NotNil(t, out)
	must.Eq(t, "channel", out.Kind)
}

func TestWithTransformer(t *testing.T) {
	tr := func(_ context.Context, _ *structpb.Struct) (*structpb.Struct, error) {
		return nil, nil
	}
	opt := WithTransformer("testname", tr)
	cfg := New(opt)
	out, ok := cfg.Transformations["testname"]
	must.True(t, ok)
	must.NotNil(t, out)
}
