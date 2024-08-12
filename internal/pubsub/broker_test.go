package pubsub

import (
	"testing"

	"github.com/shoenig/test/must"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/testutil"
)

func TestNewBroker(t *testing.T) {
	b := NewBroker("test")
	must.NotNil(t, b)
	must.Eq(t, "test", b.id)
	must.Eq(t, map[string]*Subscriber{}, b.subscribers)
}

func TestAddSubscriber(t *testing.T) {
	b := NewBroker("test")
	must.MapEmpty(t, b.subscribers)

	sub := &Subscriber{id: "test"}
	expected := map[string]*Subscriber{
		"test": sub,
	}
	b.AddSubscriber(sub)
	must.Eq(t, expected, b.subscribers)
}

func TestBroadcast(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	sub := NewSubscriber("testSub", ch)
	b := NewBroker("test")
	b.AddSubscriber(sub)

	event := testutil.BuildTestEvent()
	b.Broadcast(event)
	ev := <-ch
	must.Eq(t, event, ev)
}
