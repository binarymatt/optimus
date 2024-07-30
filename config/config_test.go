package config

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestConfigInit(t *testing.T) {
	cfg := &Config{}
	cfg.Init()
	must.Eq(t, ":8080", cfg.ListenAddress)
}
