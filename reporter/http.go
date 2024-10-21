package reporter

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hashicorp/go-retryablehttp"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var _ ReporterInterface = (*HttpReporter)(nil)

type HttpReporter struct {
	url    string
	client *retryablehttp.Client
}

func (h *HttpReporter) ReportEvent(ctx context.Context, event *optimusv1.LogEvent, status string) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	slog.Info("event reportred", "status", status)
	req, err := retryablehttp.NewRequest("POST", h.url, body)
	if err != nil {
		return err
	}
	_, err = h.client.Do(req)
	return err
}

func NewHTTPReporter(url string) *HttpReporter {
	client := retryablehttp.NewClient()
	h := &HttpReporter{
		client: client,
	}
	return h
}
