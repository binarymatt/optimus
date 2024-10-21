package destination

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestChannelSetup(t *testing.T) {
	c := &ChannelDestination{}
	must.Nil(t, c.Setup())
}

func TestChannelClose(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	c := &ChannelDestination{
		Output: ch,
	}
	must.Nil(t, c.Close())
	x, ok := <-ch
	must.Nil(t, x)
	must.False(t, ok)

}

func TestDeliverMissingChannel(t *testing.T) {
	c := &ChannelDestination{}
	err := c.Deliver(context.Background(), &optimusv1.LogEvent{})
	must.ErrorIs(t, ErrMissingChannel, err)
}

func TestDeliver(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	cd := &ChannelDestination{
		Output: ch,
	}
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		err := cd.Deliver(context.Background(), &optimusv1.LogEvent{Id: "test"})
		return err
	})
	eg.Go(func() error {
		data := <-ch
		must.Eq(t, &optimusv1.LogEvent{Id: "test"}, data)
		return nil
	})
	must.NoError(t, eg.Wait())
}
