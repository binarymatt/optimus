package input

import (
	"os"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/binarymatt/optimus/internal/pubsub"
)

func TestInitialize(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	must.NoError(t, err)
	defer os.Remove(f.Name())

	fi := &FileInput{
		Path: f.Name(),
	}
	broker := pubsub.NewBroker("test_broker")
	must.Nil(t, fi.tracker)
	err = fi.Initialize("test", broker)
	must.NoError(t, err)
	must.Eq(t, "test", fi.id)
	must.NotNil(t, fi.tracker)
}
