package transformation

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/shoenig/test/must"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/mocks"
)

func TestNew(t *testing.T) {
	mocked := mocks.NewMockTransformerImpl(t)
	mocked.EXPECT().Initialize().Return(nil)
	tr, err := New("test_name", "test", 1, []string{}, mocked)
	must.NoError(t, err)
	must.NotNil(t, tr.impl)
}

func TestProcess_HappyPath(t *testing.T) {
	mocked := mocks.NewMockTransformerImpl(t)
	ctx, cancel := context.WithCancel(context.Background())
	mocked.EXPECT().Initialize().Return(nil)
	tr, _ := New("test_name", "test", 1, []string{}, mocked)
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return tr.Process(ctx)
	})
	data, err := structpb.NewStruct(map[string]any{})
	mocked.On("Transform", ctx, data).Return(data, nil).Once()
	must.NoError(t, err)
	evt := &optimusv1.LogEvent{
		Id:   "test",
		Data: data,
	}
	eg.Go(func() error {
		tr.Subscriber.Signal(evt)
		time.Sleep(10 * time.Millisecond)
		cancel()
		return nil
	})
	err = eg.Wait()
	must.NoError(t, err)

}

func TestHclImpl(t *testing.T) {

	cases := []struct {
		name     string
		kind     string
		body     string
		expected TransformerImpl
		diags    hcl.Diagnostics
	}{
		{
			name: "jsonata transform",
			kind: KindJsonata,
			body: `expression = "$sum(orders.(price * quantity))"`,
			expected: &JsonataTransformer{
				Expression: "$sum(orders.(price * quantity))",
			},
		},
		{
			name: "jmespath transform",
			kind: KindJmespath,
			body: `expression = "locations[?state == 'WA']"`,
			expected: &JmesTransformer{
				Expression: `locations[?state == 'WA']`,
			},
		},
		{
			name: "invalid input",
			kind: "test",
			diags: append(hcl.Diagnostics{}, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "invalid transformation",
				Detail:   "test is not a valid transformation type",
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
			impl, diags := HclImpl(tc.kind, f.Body)
			must.Eq(t, tc.diags, diags)
			must.Eq(t, tc.expected, impl)
		})
	}
}
