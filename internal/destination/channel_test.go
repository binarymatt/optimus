package destination

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"

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
