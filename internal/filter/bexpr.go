package filter

import (
	"context"
	"log/slog"

	"github.com/hashicorp/go-bexpr"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/utils"
)

type BexprFilter struct {
	Expression string `yaml:"expression"`
	evaluator  *bexpr.Evaluator
}

func (b *BexprFilter) Setup() error {
	eval, err := bexpr.CreateEvaluator(b.Expression)
	if err != nil {
		slog.Error("bexpr.Setup: could not initialize evaluator", "error", err)
		return err
	}
	b.evaluator = eval
	return nil
}

func (b *BexprFilter) Process(ctx context.Context, event *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	slog.Warn("processing bexpr filter", "pattern", b.Expression)
	result, err := b.evaluator.Evaluate(event.Data.AsMap())
	if err != nil {
		slog.Error("error evaluating", "expression", b.Expression, "event_id", event.Id)
		return nil, err
	}
	if result {
		slog.Warn("passing event to next")
		return utils.CopyLogEvent(event)
	}
	return nil, nil
}
