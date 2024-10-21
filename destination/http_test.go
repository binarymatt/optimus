package destination

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestHttpSetup_Basic(t *testing.T) {
	ht := &HttpDestination{
		Endpoint: "http://example.com",
	}
	must.Nil(t, ht.client)
	err := ht.Setup()
	must.NoError(t, err)
	must.NotNil(t, ht.client)
	must.Eq(t, ht.client.RetryMax, 0)
	must.Eq(t, ht.client.HTTPClient.Timeout, 0)
}
func TestHttpSetup_Error(t *testing.T) {
	ht := &HttpDestination{}
	err := ht.Setup()
	must.ErrorIs(t, ErrMissingEndpoint, err)
}
func TestHttpSetup_Advanced(t *testing.T) {

	ht := &HttpDestination{
		Endpoint: "http://example.com",
		Retries:  2,
		Timeout:  100,
	}
	err := ht.Setup()
	must.NoError(t, err)

	must.NotNil(t, ht.client)
	must.Eq(t, 2, ht.client.RetryMax)
	must.Eq(t, 100*time.Millisecond, ht.client.HTTPClient.Timeout)
}

func TestHttpAddHeaders(t *testing.T) {
	basicEvent := &optimusv1.LogEvent{}
	cases := []struct {
		name      string
		headers   []Header
		event     *optimusv1.LogEvent
		assertion func(t *testing.T, r *retryablehttp.Request)
	}{
		{
			name: "basic",
			headers: []Header{
				{Key: "test", Value: "test"},
			},
			assertion: func(t *testing.T, r *retryablehttp.Request) {
				l := len(r.Header)
				must.Eq(t, 1, l)
				must.Eq(t, "test", r.Header.Get("test"))
			},
		},
		{
			name: "path",
			event: &optimusv1.LogEvent{
				Id: "one",
			},
			headers: []Header{
				{Key: "test", Value: "test", Path: "/Id"},
			},
			assertion: func(t *testing.T, r *retryablehttp.Request) {
				l := len(r.Header)
				must.Eq(t, 1, l)
				must.Eq(t, "one", r.Header.Get("test"))
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var event *optimusv1.LogEvent
			ht := &HttpDestination{
				Headers: tc.headers,
			}
			req, err := retryablehttp.NewRequest("POST", "http://example.com", nil)
			must.NoError(t, err)
			if tc.event != nil {
				event = tc.event
			} else {
				event = basicEvent
			}
			ht.AddHeaders(req, event)
			tc.assertion(t, req)
		})
	}
}

func TestHttpDeliver(t *testing.T) {
	expected := &optimusv1.LogEvent{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"sample": structpb.NewStringValue("value"),
			},
		},
	}
	expectedBody, err := json.Marshal(expected.Data.AsMap())
	must.NoError(t, err)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		must.NoError(t, err)
		must.Eq(t, expectedBody, data)
	}))
	defer svr.Close()
	ht := &HttpDestination{
		Endpoint: svr.URL,
	}
	err = ht.Setup()
	must.NoError(t, err)
	err = ht.Deliver(context.Background(), expected)
	must.NoError(t, err)
}
