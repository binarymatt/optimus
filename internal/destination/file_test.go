package destination

import (
	"bytes"
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestFileDeliver(t *testing.T) {
	w := bytes.NewBuffer(nil)
	fd := &FileDestination{
		w: w,
	}
	event := &optimusv1.LogEvent{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"test": structpb.NewStringValue("val"),
			},
		},
	}
	err := fd.Deliver(context.Background(), event)
	must.NoError(t, err)
	must.Eq(t, "{\"test\":\"val\"}\n", w.String())
}
