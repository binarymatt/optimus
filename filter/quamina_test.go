package filter

import (
	"context"
	"errors"
	"testing"

	"github.com/shoenig/test/must"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/testutil"
)

func TestQuaminaSetup(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		err     error
	}{
		{
			name:    "happy path",
			pattern: `{"Image": {"Width": [800]}}`,
		},
		{
			name:    "malformed pattern",
			pattern: `{"foo": 1}`,
			err:     errors.New("pattern malformed, illegal 1"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			qf := &QuaminaFilter{
				Patterns: map[string]string{"test": tc.pattern},
			}
			err := qf.Setup()
			must.Eq(t, tc.err, err)
		})
	}

}

func TestQuaminaProcess(t *testing.T) {
	event := testutil.BuildTestEvent()
	cases := []struct {
		name     string
		pattern  string
		expected *optimusv1.LogEvent
		err      error
	}{
		{
			name:     "happy path",
			pattern:  `{"Image": {"Width": [800]}}`,
			expected: event,
		},
		{
			name:    "no match",
			pattern: `{"Image":{"Width":[1]}}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			qf := &QuaminaFilter{
				Patterns: map[string]string{"test": tc.pattern},
			}
			err := qf.Setup()
			must.Eq(t, tc.err, err)

			res, err := qf.Process(ctx, event)
			must.Eq(t, tc.err, err)
			must.Eq(t, tc.expected, res)
		})
	}
}
