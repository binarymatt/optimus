package config

import (
	"log/slog"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/binarymatt/optimus/mocks"
)

func TestWithLogLevel(t *testing.T) {
	cases := []struct {
		name          string
		expectedLevel slog.Level
		level         string
	}{
		{
			name:          "empty string",
			expectedLevel: slog.LevelInfo,
			level:         "",
		},
		{
			name:          "debug",
			expectedLevel: slog.LevelDebug,
			level:         "debug",
		},
		{
			name:          "info",
			expectedLevel: slog.LevelInfo,
			level:         "",
		},
		{
			name:          "warn",
			expectedLevel: slog.LevelWarn,
			level:         "warn",
		},
		{
			name:          "error",
			expectedLevel: slog.LevelError,
			level:         "error",
		},
		{
			name:          "upper case string",
			expectedLevel: slog.LevelInfo,
			level:         "INFO",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := New(WithLogLevel(tc.level))
			must.Eq(t, tc.expectedLevel, cfg.LogLevel)
		})
	}
}
func skipInitialize() ConfigOption {
	return func(c *Config) {
		c.initialized = true
	}
}

func TestWithListenAddress(t *testing.T) {
	cfg := New(skipInitialize())
	WithListenAddress("localhost:9090")(cfg)
	must.Eq(t, "localhost:9090", cfg.ListenAddress)
}

func TestWithMetricsEnabled(t *testing.T) {
	cfg := New(skipInitialize())
	WithMetricsEnabled()(cfg)
	must.True(t, cfg.MetricsEnabled)
	must.True(t, cfg.HttpInputEnabled)
}

func TestWithTransformer(t *testing.T) {
	trImpl := mocks.NewMockTransformerImpl(t)
	trImpl.EXPECT().Initialize().Return(nil)
	opt := WithTransformer("testname", "test", 1, []string{}, trImpl)
	cfg := New(opt)
	out := cfg.Transformations[0]
	must.NotNil(t, out)
}
