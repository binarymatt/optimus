package pubsub

import (
	"sync"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type Subscriber interface {
	Destruct()
	Signal(*optimusv1.LogEvent)
	GetID() string
}

var _ Subscriber = (*subscriber)(nil)

type subscriber struct {
	id       string                   // id of subscriber
	messages chan *optimusv1.LogEvent // messages channel
	//topics   map[string]bool          // topics it is subscribed to.
	active bool         // if given subscriber is active
	mutex  sync.RWMutex // lock
}

func NewSubscriber(id string, messages chan *optimusv1.LogEvent) *subscriber {
	return &subscriber{
		id:       id,
		messages: messages,
		active:   true,
	}
}
func (s *subscriber) Destruct() {
	// destructor for subscriber.
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.active = false
	close(s.messages)
}
func (s *subscriber) Signal(le *optimusv1.LogEvent) {
	// Gets the message from the channel
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.messages <- le
}
func (s *subscriber) GetID() string {
	return s.id
}
