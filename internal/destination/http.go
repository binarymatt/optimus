package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var (
	ErrMissingEndpoint = errors.New("missing endpoint information")
)

type HttpDestination struct {
	endpoint string
	method   string
	client   *retryablehttp.Client
}

func (h *HttpDestination) Setup(cfg map[string]any) error {
	timeout := 0
	retries := 0
	method := "POST"
	endpoint := ""
	if val, ok := cfg["timeout"]; ok {
		timeout = val.(int)
	}
	if val, ok := cfg["retries"]; ok {
		retries = val.(int)
	}

	if val, ok := cfg["method"]; ok {
		method = val.(string)
	}
	if val, ok := cfg["endpoint"]; ok {
		endpoint = val.(string)
	}
	if endpoint == "" {
		return ErrMissingEndpoint
	}
	c := retryablehttp.NewClient()
	c.HTTPClient.Timeout = time.Duration(timeout) * time.Millisecond
	c.RetryMax = retries
	h.client = c
	h.method = method
	return nil
}
func (h *HttpDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	raw, err := json.Marshal(event.Data.AsMap())
	if err != nil {
		return err
	}
	req, err := retryablehttp.NewRequest(h.method, h.endpoint, bytes.NewBuffer(raw))
	if err != nil {
		return err
	}
	resp, err := h.client.Do(req)
	slog.Debug("http delivery done", "code", resp.StatusCode)
	return err
}
