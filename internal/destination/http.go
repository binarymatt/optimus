package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/mitchellh/pointerstructure"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

var (
	ErrMissingEndpoint = errors.New("missing endpoint information")
)

type Header struct {
	Key   string `hcl:"key"`
	Value string `hcl:"value"`
	Path  string `hcl:"path"`
}
type HttpDestination struct {
	client   *retryablehttp.Client
	Timeout  int      `hcl:"timeout"`
	Retries  int      `hcl:"retries"`
	Endpoint string   `hcl:"endpoint"`
	Method   string   `hcl:"http_method"`
	Headers  []Header `hcl:"headers"`
}

func (h *HttpDestination) Setup() error {

	if h.Endpoint == "" {
		return ErrMissingEndpoint
	}
	c := retryablehttp.NewClient()
	c.HTTPClient.Timeout = time.Duration(h.Timeout) * time.Millisecond
	c.RetryMax = h.Retries

	h.client = c
	return nil
}
func (h *HttpDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	raw, err := json.Marshal(event.Data.AsMap())
	if err != nil {
		return err
	}
	req, err := retryablehttp.NewRequest(h.Method, h.Endpoint, bytes.NewBuffer(raw))
	if err != nil {
		return err
	}
	h.AddHeaders(req, event)

	resp, err := h.client.Do(req)
	slog.Debug("http delivery done", "code", resp.StatusCode)
	return err
}

func (h *HttpDestination) AddHeaders(req *retryablehttp.Request, event *optimusv1.LogEvent) {
	for _, header := range h.Headers {
		key := header.Key
		val := header.Value
		if header.Path != "" {
			var err error
			pval, err := pointerstructure.Get(event, header.Path)
			if err != nil {
				slog.Error("could not get header value via path", "error", err, "key", key)
			}
			if pval != "" {
				val = fmt.Sprintf("%v", pval)
			}
		}
		req.Header.Add(key, val)

	}
}

func (h *HttpDestination) Close() error {
	return nil
}
