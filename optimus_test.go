package optimus

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/binarymatt/optimus/config"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	o := New(cfg)
	must.NotNil(t, o)
}
