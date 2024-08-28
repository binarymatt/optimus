package pubsub

import (
	"testing"

	"github.com/shoenig/test/must"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestNewSubscriber(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	s := NewSubscriber("testid", ch)
	must.Eq(t, "testid", s.id)
	must.Eq(t, ch, s.messages)
	must.True(t, s.active)
}

func TestDestruct(t *testing.T) {
	ch := make(chan *optimusv1.LogEvent)
	s := NewSubscriber("testid", ch)
	s.Destruct()
	must.False(t, s.active)
	_, ok := <-ch
	must.False(t, ok)
}
