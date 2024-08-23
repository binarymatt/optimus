package utils

import (
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/testutil"
)

func TestCopyLogEvent(t *testing.T) {
	payload, err := structpb.NewStruct(map[string]any{
		"hello": "world",
	})
	must.NoError(t, err)
	event := &optimusv1.LogEvent{
		Id:        "test_event",
		Data:      payload,
		Source:    "test_input",
		Upstreams: []string{"test"},
	}
	newEvent, err := CopyLogEvent(event)
	must.NoError(t, err)
	must.Eq(t, event, newEvent, testutil.CmpTransform)
}

func TestCopyLogEvent_Nil(t *testing.T) {
	newEvent, err := CopyLogEvent(nil)
	must.ErrorIs(t, err, ErrNilEvent)
	must.Nil(t, newEvent)
}
