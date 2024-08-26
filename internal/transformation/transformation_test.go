package transformation

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestNew(t *testing.T) {

	tr := New("test_name", nil)
	must.Eq(t, "test_name", tr.Name)
	must.Nil(t, tr.transformer)
	transform := func(ctx context.Context, _ *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
		return nil, nil
	}
	tr = New("test_name", transform)
	must.NotNil(t, tr.transformer)
}
