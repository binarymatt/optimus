package destination

import (
	"bytes"
	"context"
	"testing"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func TestStdSetup(t *testing.T) {
	sd := &StdOutDestination{}
	must.Nil(t, sd.encoder)
	err := sd.Setup()
	must.NoError(t, err)
	must.NotNil(t, sd.encoder)
}

func TestStdOutDeliver(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	sd := StdOutDestination{
		Writer: writer,
	}
	err := sd.Setup()
	must.NoError(t, err)

	event := &optimusv1.LogEvent{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"test": structpb.NewStringValue("val"),
			},
		},
	}
	err = sd.Deliver(context.Background(), event)
	must.NoError(t, err)
	must.Eq(t, `{"test":"val"}`, writer.String())
}
