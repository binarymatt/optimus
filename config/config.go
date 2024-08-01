package config

import (
	"os"

	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
)

type Config struct {
	DataDir          string `yaml:"data_dir"` // used to store data about positions.
	MetricsEnabled   bool   `yaml:"metrics_enabled"`
	LogLevel         string `yaml:"log_level"`
	HttpInputEnabled bool
	ListenAddress    string                              `yaml:"listen_address"`
	Inputs           map[string]*input.Input             `yaml:"inputs"`
	Filters          map[string]*filter.Filter           `yaml:"filters"`
	Destinations     map[string]*destination.Destination `yaml:"destinations"`
}

func LoadYamlFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadYaml(data)
}

func LoadYaml(data []byte) (*Config, error) {
	cfg := Config{}
	err := yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func (c *Config) Init() {
	if c.ListenAddress == "" {
		c.ListenAddress = ":8080"
	}
}

func (c *Config) WithChannelInput(name string, in <-chan *optimusv1.LogEvent) *Config {
	ci := &input.ChannelInput{
		Input: in,
	}
	i := &input.Input{
		ID:   name,
		Kind: "channel",
	}
	i.WithInputProcessor(ci)
	c.Inputs[name] = i
	return c
}

// NOTE update to include channel output implementation
func (c *Config) WithChannelOutput(name string, out chan<- *optimusv1.LogEvent) *Config {
	return c
}
