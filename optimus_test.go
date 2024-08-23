package optimus

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/binarymatt/optimus/config"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	o, err := New(cfg)
	must.NoError(t, err)
	must.NotNil(t, o)
}
