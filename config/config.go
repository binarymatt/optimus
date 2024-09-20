package config

import (
	"os"

	"gopkg.in/yaml.v3"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
	"github.com/binarymatt/optimus/internal/transformation"
)

type ConfigOption func(*Config)
type Config struct {
	DataDir          string `yaml:"data_dir"` // used to store data about positions.
	MetricsEnabled   bool   `yaml:"metrics_enabled"`
	LogLevel         string `yaml:"log_level"`
	HttpInputEnabled bool
	ListenAddress    string                                    `yaml:"listen_address"`
	Inputs           map[string]*input.Input                   `yaml:"inputs"`
	Filters          map[string]*filter.Filter                 `yaml:"filters"`
	Destinations     map[string]*destination.Destination       `yaml:"destinations"`
	Transformations  map[string]*transformation.Transformation `yaml:"transformations"`
}

func LoadYamlFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadYaml(data)
}

func LoadYaml(data []byte) (*Config, error) {
	cfg := Config{
		Inputs:          make(map[string]*input.Input),
		Filters:         make(map[string]*filter.Filter),
		Destinations:    make(map[string]*destination.Destination),
		Transformations: make(map[string]*transformation.Transformation),
	}
	err := yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func (c *Config) Init() {
	if c.ListenAddress == "" {
		c.ListenAddress = ":8080"
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
	}
}

func WithInput(id, kind string, impl input.InputProcessor) ConfigOption {
	return func(c *Config) {
		c.Inputs[id] = input.New(id, kind, impl)
		if kind == "http" {
			c.HttpInputEnabled = true
		}
	}
}
func WithDestination(id, kind string, subscriptions []string, impl destination.DestinationProcessor) ConfigOption {
	return func(c *Config) {
		c.Destinations[id] = destination.New(id, kind, subscriptions, impl)
	}
}
func WithFilter(id, kind string, subscriptions []string, impl filter.FilterProcessor) ConfigOption {
	return func(c *Config) {
		c.Filters[id] = filter.New(id, kind, subscriptions, impl)
	}
}

func WithChannelInput(name string, in <-chan *optimusv1.LogEvent) ConfigOption {
	return func(c *Config) {
		ci := &input.ChannelInput{
			Input: in,
		}
		c.Inputs[name] = input.New(name, "channel", ci)
	}
}

func WithChannelOutput(name string, out chan<- *optimusv1.LogEvent, subscriptions []string) ConfigOption {
	return func(c *Config) {
		cd := &destination.ChannelDestination{
			Output: out,
		}
		c.Destinations[name] = destination.New(name, "channel", subscriptions, cd)
	}
}

func WithTransformer(id, kind string, subscriptions []string, transformer transformation.TransformerImpl) ConfigOption {
	return func(c *Config) {
		t := transformation.New(id, kind, subscriptions, transformer)
		c.Transformations[id] = t
	}
}
func NewWithYaml(yaml []byte, opts ...ConfigOption) (*Config, error) {
	cfg, err := LoadYaml(yaml)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg, nil
}
func New(opts ...ConfigOption) *Config {
	c := &Config{
		Inputs:          make(map[string]*input.Input),
		Filters:         make(map[string]*filter.Filter),
		Destinations:    make(map[string]*destination.Destination),
		Transformations: make(map[string]*transformation.Transformation),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
