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
func TestLoadYaml(t *testing.T) {
	str := `---
data_dir: "/tmp/data"
metrics_enabled: true
listen_address: ":8080"
console: true
log_level: debug
inputs:
  fileInput:
    kind: file
    path: "./cmd/test/tmp"
  httpInput:
    kind: http
destinations:
  sampleout:
    kind: stdout
    subscriptions:
      - fileInput
      - httpInput
      - testing
  samplefile:
    kind: file
    path: "test.ndjson"
    subscriptions:
      - httpInput
`
	cfg, err := LoadYaml([]byte(str))
	must.NoError(t, err)
	must.Eq(t, "stdout", cfg.Destinations["sampleout"].Kind)
	must.Eq(t, "file", cfg.Destinations["samplefile"].Kind)
	//must.Eq(t, cfg.Destinations, expected)
}
