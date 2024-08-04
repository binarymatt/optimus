package destination

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type TestProcessor struct {
	SetupCalled    bool
	DeliverCalled  bool
	CloseCalled    bool
	ProcessedEvent *optimusv1.LogEvent
}

func (tp *TestProcessor) Setup() error {
	tp.SetupCalled = true
	return nil
}
func (tp *TestProcessor) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	fmt.Println("delivering")
	tp.DeliverCalled = true
	tp.ProcessedEvent = event
	return nil
}
func (tp *TestProcessor) Close() error {
	tp.CloseCalled = true
	return nil
}
func TestInit(t *testing.T) {
	p := &TestProcessor{}
	d := &Destination{}
	d.WithProcessor(p)
	must.Eq(t, "", d.id)
	must.Eq(t, 0, d.BufferSize)
	must.Nil(t, d.inputs)
	must.Nil(t, d.Subscriber)

	err := d.Init("testing")
	must.NoError(t, err)
	must.Eq(t, "testing", d.id)
	must.True(t, p.SetupCalled)
}

func TestInit_Error(t *testing.T) {
	testErr := errors.New("oops")
	d := &Destination{
		initialize: func() error {
			return testErr
		},
	}
	err := d.Init("testing")
	must.ErrorIs(t, testErr, err)
}

func TestWithProcessor(t *testing.T) {
	d := &Destination{}
	must.Nil(t, d.initialize)
	must.Nil(t, d.process)
	must.Nil(t, d.closer)
	d.WithProcessor(&TestProcessor{})
	must.NotNil(t, d.initialize)
	must.NotNil(t, d.process)
	must.NotNil(t, d.closer)
}

func TestProcess(t *testing.T) {
	d := &Destination{}
	p := &TestProcessor{}
	d.WithProcessor(p)
	err := d.Init("testing")
	must.NoError(t, err)
	event := &optimusv1.LogEvent{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"test": structpb.NewStringValue("val"),
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		d.Process(ctx)
		return nil
	})
	eg.Go(func() error {
		d.inputs <- event
		time.Sleep(5 * time.Millisecond)
		cancel()
		return nil
	})
	err = eg.Wait()
	must.NoError(t, err)
	must.True(t, p.DeliverCalled)
	must.Eq(t, event, p.ProcessedEvent)
}

var data = `---
kind: stdout
subscriptions:
  - fileInput
  - httpInput
  - testing
`

var dataErr = `---
kind: unknown
subscriptions:
  - fileInput
  - httpInput
  - testing
`

func TestUnmarshalYaml(t *testing.T) {
	var raw yaml.Node
	d := &Destination{}
	err := yaml.Unmarshal([]byte(data), &raw)
	must.NoError(t, err)
	err = d.UnmarshalYAML(&raw)
	must.NoError(t, err)
}
func TestUnmarshalYaml_NoInternal(t *testing.T) {
	var raw yaml.Node
	d := &Destination{}
	err := yaml.Unmarshal([]byte(dataErr), &raw)
	must.NoError(t, err)
	err = d.UnmarshalYAML(&raw)
	must.ErrorIs(t, ErrNoProcessor, err)
}
