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
	Timeout  int    `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"http_method"`
}

func (h *HttpDestination) Setup() error {

	if h.Endpoint == "" {
		return ErrMissingEndpoint
	}
	c := retryablehttp.NewClient()
	c.HTTPClient.Timeout = time.Duration(h.Timeout) * time.Millisecond
	c.RetryMax = h.Retries

	h.method = h.Method
	h.client = c
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
