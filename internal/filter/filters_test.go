package filter

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/mocks"
)

func TestInit(t *testing.T) {
	mocked := mocks.NewMockFilterProcessor(t)
	mocked.EXPECT().Setup().Return(nil)
	f := &Filter{
		ID:   "test",
		impl: mocked,
	}
	must.Nil(t, f.Broker)
	must.Nil(t, f.Subscriber)
	must.Nil(t, f.inputs)
	must.Zero(t, f.BufferSize)

	_, err := f.Init()
	must.NoError(t, err)
	must.Eq(t, 5, f.BufferSize)
	must.NotNil(t, f.Broker)
	must.NotNil(t, f.inputs)
	must.NotNil(t, f.Subscriber)
}

func TestProcess(t *testing.T) {
	event := &optimusv1.LogEvent{
		Id: "test",
	}
	mocked := mocks.NewMockFilterProcessor(t)
	mocked.EXPECT().Setup().Return(nil)
	f := &Filter{
		BufferSize: 1,
		impl:       mocked,
	}
	_, err := f.Init()
	must.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mocked.EXPECT().Process(ctx, event).Return(event, nil).Once()
	eg.Go(func() error {
		// send event to channel
		f.inputs <- event
		// cancel context
		time.Sleep(2 * time.Millisecond)
		cancel()
		return nil
	})
	eg.Go(func() error {
		return f.Process(ctx)
	})
	must.NoError(t, eg.Wait())

}
func TestHclImpl(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		body     string
		expected FilterProcessor
		diags    hcl.Diagnostics
	}{
		{
			name: "bexpr",
			kind: KindBexpr,
			body: `
			expression = "test == true"
			`,
			expected: &BexprFilter{
				Expression: `test == true`,
			},
		},
		{
			name: "quamina",
			kind: KindQuamina,
			body: `patterns = {
				test = "true"
			}`,
			expected: &QuaminaFilter{
				Patterns: map[string]string{
					"test": "true",
				},
			},
		},
		{
			name: "invalid",
			kind: "test",
			diags: append(hcl.Diagnostics{}, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "invalid filter",
				Detail:   "test is not a valid filter type.",
			}),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, diag := hclparse.NewParser().ParseHCL([]byte(tc.body), "test.file")
			if diag.HasErrors() {
				t.Logf("parse errors: %v", diag.Errs())
			}
			must.False(t, diag.HasErrors())
			impl, diags := HclImpl(tc.kind, nil, f.Body)
			must.Eq(t, tc.diags, diags)
			must.Eq(t, tc.expected, impl)

		})
	}
}
