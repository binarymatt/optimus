package transformation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestNew(t *testing.T) {

	tr := New("test_name", nil)
	must.Eq(t, "test_name", tr.Name)
	must.Nil(t, tr.transformer)
	transform := func(ctx context.Context, _ *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
		return nil, nil
	}
	tr = New("test_name", transform)
	must.NotNil(t, tr.transformer)
}

func TestProcess_HappyPath(t *testing.T) {
	transformed := false
	transform := func(ctx context.Context, evt *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
		fmt.Println("transoforming")
		transformed = true
		return evt, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	tr := New("test_name", transform)
	_ = tr.Init()
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return tr.Process(ctx)
	})
	evt := &optimusv1.LogEvent{
		Id: "test",
	}
	eg.Go(func() error {
		tr.Subscriber.Signal(evt)
		time.Sleep(10 * time.Millisecond)
		cancel()
		return nil
	})
	err := eg.Wait()
	must.NoError(t, err)
	must.True(t, transformed)

}
