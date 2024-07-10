package filter

import (
	"context"
	"encoding/json"

	"quamina.net/go/quamina"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type QuaminaFilter struct {
	Patterns map[string]string `yaml:"patterns"`
	q        *quamina.Quamina
}

func (q *QuaminaFilter) Setup() error {
	qu, err := quamina.New(quamina.WithMediaType("application/json"))
	if err != nil {
		return err
	}
	for key, value := range q.Patterns {
		if err := qu.AddPattern(key, value); err != nil {
			return err
		}
	}
	q.q = qu

	return nil
}
func (q *QuaminaFilter) Process(ctx context.Context, event *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	raw, err := json.Marshal(event.Data.AsMap())
	if err != nil {
		return nil, err
	}
	matches, err := q.q.MatchesForEvent(raw)
	if err != nil {
		return nil, err
	}
	if len(matches) > 0 {
		newEvent := &optimusv1.LogEvent{
			Id:        event.Id,
			Data:      event.Data,
			Source:    event.Source,
			Upstreams: event.Upstreams,
		}
		return newEvent, nil
	}
	return nil, nil
}
