package reporter

import (
	"context"
	"io"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var _ ReporterInterface = (*FileReporter)(nil)

type FileReporter struct {
	w io.Writer
}

func (f *FileReporter) ReportEvent(ctx context.Context, event *optimusv1.LogEvent, status string) error {
	return nil
}
