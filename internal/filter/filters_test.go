package filter

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var unknownYaml = `---
kind: test
subscriptions:
  - fileInput
BufferSize: 2
`

var yamlStr = `---
kind: bexpr
subscriptions:
  - fileInput
BufferSize: 1
expression: test.foo == true
`

func TestInit(t *testing.T) {
	f := &Filter{}
	must.Nil(t, f.Broker)
	must.Nil(t, f.Subscriber)
	must.Nil(t, f.inputs)
	must.Nil(t, f.process)
	must.Zero(t, f.BufferSize)
	must.Eq(t, "", f.id)

	f.Init("testing")
	must.Eq(t, "testing", f.id)
	must.Eq(t, 5, f.BufferSize)
	must.NotNil(t, f.Broker)
	must.NotNil(t, f.inputs)
	must.NotNil(t, f.Subscriber)
}

func TestSetupInternal(t *testing.T) {
	f := &Filter{}
	must.NoError(t, f.SetupInternal())
}

func TestUnknownFilter(t *testing.T) {
	var raw yaml.Node
	f := &Filter{}
	must.NoError(t, yaml.Unmarshal([]byte(unknownYaml), &raw))
	must.NoError(t, f.UnmarshalYAML(&raw))
}

func TestUnmarshalYaml(t *testing.T) {
	var raw yaml.Node
	f := &Filter{}
	must.NoError(t, yaml.Unmarshal([]byte(yamlStr), &raw))
	must.NoError(t, f.UnmarshalYAML(&raw))
}

func TestProcess(t *testing.T) {
	f := &Filter{
		BufferSize: 1,
		process: func(ctx context.Context, event *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
			return event, nil
		},
	}
	f.Init("testing")
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// send event to channel
		f.inputs <- &optimusv1.LogEvent{
			Id: "test",
		}
		// cancel context
		cancel()
		return nil
	})
	eg.Go(func() error {
		return f.Process(ctx)
	})
	must.NoError(t, eg.Wait())

}
