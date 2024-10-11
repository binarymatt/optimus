package config

import (
	"log/slog"
	"strings"

	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
	"github.com/binarymatt/optimus/internal/transformation"
)

type ConfigOption func(*Config)

func WithLogLevel(logLevel string) ConfigOption {
	return func(c *Config) {
		level := slog.LevelInfo
		switch strings.ToLower(logLevel) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
		// slog.SetLogLoggerLevel(level)
		c.LogLevel = level
	}
}
func WithListenAddress(address string) ConfigOption {
	return func(c *Config) {
		c.ListenAddress = address
		c.HttpInputEnabled = true
	}
}

func WithMetricsEnabled() ConfigOption {
	return func(c *Config) {
		c.MetricsEnabled = true
		c.HttpInputEnabled = true
	}
}

func WithInput(id, kind string, impl input.InputProcessor) ConfigOption {
	return func(c *Config) {
		in, err := input.New(id, kind, impl)
		if err != nil {
			slog.Error("error creating input", "error", err, "input", in, "impl", impl)
			return
		}
		c.Inputs = append(c.Inputs, in)
		c.references[id] = in.Broker

	}
}

func WithDestination(id, kind string, bufferSize int, subscriptions []string, impl destination.DestinationProcessor) ConfigOption {
	return func(c *Config) {
		d, err := destination.New(id, kind, bufferSize, subscriptions, impl)
		if err != nil {
			slog.Error("error creating destination", "error", err)
			return
		}
		c.Destinations = append(c.Destinations, d)
	}
}
func WithFilter(id, kind string, bufferSize int, subscriptions []string, impl filter.FilterProcessor) ConfigOption {
	return func(c *Config) {
		f, err := filter.New(id, kind, bufferSize, subscriptions, impl)
		if err != nil {
			slog.Error("error creating filter", "error", err)
			return
		}
		c.Filters = append(c.Filters, f)
		c.references[id] = f.Broker
	}
}

func WithTransformer(id, kind string, bufferSize int, subscriptions []string, transformer transformation.TransformerImpl) ConfigOption {
	return func(c *Config) {
		t, err := transformation.New(id, kind, bufferSize, subscriptions, transformer)
		if err != nil {
			slog.Error("coudl not create new transformation", "error", err)
			return
		}
		c.Transformations = append(c.Transformations, t)
		c.references[id] = t.Broker
	}
}
