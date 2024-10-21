package reporter

import (
	"context"
	"database/sql"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var _ ReporterInterface = (*SqliteReporter)(nil)

type SqliteReporter struct {
	db sql.DB
}

func (s *SqliteReporter) ReportEvent(ctx context.Context, event *optimusv1.LogEvent, status string) error {
	return nil
}
