package input

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/binarymatt/optimus/mocks"
)

func TestHttpInitialize(t *testing.T) {
	hi := &HTTPInput{}
	broker := mocks.NewMockBroker(t)
	must.NoError(t, hi.Initialize("test", broker))
	must.Eq(t, "test", hi.ID)
}
