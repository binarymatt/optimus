package reporter

import (
	"context"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type ReporterInterface interface {
	ReportEvent(ctx context.Context, event *optimusv1.LogEvent, status string) error
}
