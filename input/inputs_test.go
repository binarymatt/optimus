package input

import (
	"errors"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"github.com/binarymatt/optimus/mocks"
)

func TestNew(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	mocked.EXPECT().Initialize("test", mock.AnythingOfType("*pubsub.broker")).Return(nil)
	p, err := New("test", KindHttp, mocked)
	must.NoError(t, err)
	must.NotNil(t, p)
}
func TestProcess(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		impl: mocked,
	}
	ctx := context.Background()
	mocked.EXPECT().Process(ctx).Return(nil).Once()
	must.NoError(t, i.Process(ctx))
}

func TestProcess_Error(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		impl: mocked,
	}
	ctx := context.Background()
	errOops := errors.New("oops")
	mocked.EXPECT().Process(ctx).Return(errOops).Once()
	must.ErrorIs(t, errOops, i.Process(ctx))
}

func TestInit(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		ID:   "testid",
		impl: mocked,
	}
	mocked.EXPECT().
		Initialize("testid", mock.AnythingOfType("*pubsub.broker")).
		Return(nil).Once()
	err := i.Init()
	must.NotNil(t, i.Broker)
	must.NoError(t, err)
}

func TestInit_Error(t *testing.T) {
	mocked := mocks.NewMockInputProcessor(t)
	i := &Input{
		ID:   "testid",
		impl: mocked,
	}
	errOops := errors.New("oops")
	mocked.EXPECT().
		Initialize("testid", mock.AnythingOfType("*pubsub.broker")).
		Return(errOops).Once()
	err := i.Init()
	must.Nil(t, i.Broker)
	must.ErrorIs(t, errOops, err)
}

func TestHclImpl(t *testing.T) {

	cases := []struct {
		name     string
		kind     string
		body     string
		expected InputProcessor
		diags    hcl.Diagnostics
	}{
		{
			name: "http input",
			kind: KindFile,
			body: `path = "test.out"`,
			expected: &FileInput{
				Path: "test.out",
			},
		},
		{
			name:     "http input",
			kind:     KindHttp,
			expected: &HTTPInput{},
		},
		{
			name: "invalid input",
			kind: "test",
			diags: append(hcl.Diagnostics{}, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "invalid input",
				Detail:   "test is not a valid input type",
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
