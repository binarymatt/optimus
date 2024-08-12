package pubsub

import (
	"log/slog"
	"sync"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type Subscribers map[string]*Subscriber
type Broker struct {
	id          string
	subscribers Subscribers  // map of subscribers id:Subscriber
	mut         sync.RWMutex // mutex lock
}

func NewBroker(id string) *Broker {
	// returns new broker object
	return &Broker{
		id:          id,
		subscribers: map[string]*Subscriber{},
	}
}
func (b *Broker) AddSubscriber(s *Subscriber) {
	// Add subscriber to the broker.
	b.mut.Lock()
	defer b.mut.Unlock()
	b.subscribers[s.id] = s
}
func (b *Broker) Broadcast(event *optimusv1.LogEvent) {
	slog.Debug("broadcasting...", "subscribers", len(b.subscribers))
	// broadcast message to all topics.
	event.Upstreams = append(event.Upstreams, b.id)
	for _, s := range b.subscribers {
		slog.Debug("sending to subscriber", "id", s.id)
		go func(s *Subscriber) {
			s.Signal(event)
		}(s)

	}
}
