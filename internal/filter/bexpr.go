package filter

import (
	"context"
	"log/slog"

	"github.com/hashicorp/go-bexpr"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type BexprFilter struct {
	Expression string `yaml:"expression"`
	evaluator  *bexpr.Evaluator
}

func (b *BexprFilter) Setup() error {
	eval, err := bexpr.CreateEvaluator(b.Expression)
	if err != nil {
		slog.Error("could not initialize evaluator")
		return err
	}
	b.evaluator = eval
	return nil
}

func (b *BexprFilter) Process(ctx context.Context, event *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	result, err := b.evaluator.Evaluate(event.Data.AsMap())
	if err != nil {
		slog.Error("error evaluating", "expression", b.Expression, "event_id", event.Id)
		return nil, err
	}
	if result {
		new := &optimusv1.LogEvent{
			Id:        event.Id,
			Data:      event.Data,
			Source:    event.Source,
			Upstreams: event.Upstreams,
		}
		return new, nil
	}
	return nil, nil
}
