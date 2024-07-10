package config

import (
	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
)

type Config struct {
	DataDir        string                              `yaml:"data_dir"` // used to store data about positions.
	MetricsEnabled bool                                `yaml:"metrics_enabled"`
	ListenAddress  string                              `yaml:"listen_address"`
	Inputs         map[string]*input.Input             `yaml:"inputs"`
	Filters        map[string]*filter.Filter           `yaml:"filters"`
	Destinations   map[string]*destination.Destination `yaml:"destinations"`
}
