package pubsub

import (
	"log/slog"
	"sync"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type Subscribers map[string]Subscriber
type Broker interface {
	AddSubscriber(Subscriber)
	Broadcast(*optimusv1.LogEvent)
}
type broker struct {
	id          string
	subscribers Subscribers  // map of subscribers id:Subscriber
	mut         sync.RWMutex // mutex lock
}

var _ Broker = (*broker)(nil)

func NewBroker(id string) *broker {
	// returns new broker object
	return &broker{
		id:          id,
		subscribers: map[string]Subscriber{},
	}
}
func (b *broker) AddSubscriber(s Subscriber) {
	// Add subscriber to the broker.
	b.mut.Lock()
	defer b.mut.Unlock()
	b.subscribers[s.GetID()] = s
}
func (b *broker) Broadcast(event *optimusv1.LogEvent) {
	slog.Debug("broadcasting...", "subscribers", len(b.subscribers))
	// broadcast message to all topics.
	event.Upstreams = append(event.Upstreams, b.id)
	for _, s := range b.subscribers {
		slog.Info("sending to subscriber", "id", s.GetID())
		go func(s Subscriber) {
			s.Signal(event)
		}(s)

	}
}
