package reporter

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var _ ReporterInterface = (*HttpReporter)(nil)

type HttpReporter struct {
	url    string
	client *retryablehttp.Client
}

func (h *HttpReporter) ReportEvent(ctx context.Context, event *optimusv1.LogEvent, string status) error {
	return nil
}

func NewHTTPReporter(url string) {
	client := retryablehttp.NewClient()
	h := &HttpReporter{
		client: client,
	}
}
