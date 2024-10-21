package destination

import (
	"bytes"
	"context"
	"os"
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

func TestFileSetup_MissingPath(t *testing.T) {

	fd := &FileDestination{}
	err := fd.Setup()
	must.ErrorIs(t, ErrMissingPath, err)
}

func TestFileSetup(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	defer os.Remove(f.Name())

	must.NoError(t, err)
	fd := &FileDestination{
		Path: f.Name(),
	}
	must.Nil(t, fd.f)
	must.Nil(t, fd.w)
	err = fd.Setup()
	must.NoError(t, err)
	must.NotNil(t, fd.f)
	must.NotNil(t, fd.w)
}

func TestFileClose(t *testing.T) {

	f, err := os.CreateTemp("", "sample")
	defer os.Remove(f.Name())

	must.NoError(t, err)
	fd := &FileDestination{
		Path: f.Name(),
	}
	err = fd.Setup()
	must.NoError(t, err)
	must.NoError(t, fd.Close())
	must.ErrorIs(t, fd.Close(), os.ErrClosed)
}
