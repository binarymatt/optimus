package config

import (
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"

	"github.com/binarymatt/optimus/internal/destination"
	"github.com/binarymatt/optimus/internal/filter"
	"github.com/binarymatt/optimus/internal/input"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/transformation"
)

type Config struct {
	// DataDir          string // used to store data about positions.
	MetricsEnabled   bool
	LogLevel         slog.Level
	HttpInputEnabled bool
	ListenAddress    string
	references       map[string]pubsub.Broker
	Inputs           []*input.Input
	Filters          []*filter.Filter
	Destinations     []*destination.Destination
	Transformations  []*transformation.Transformation
}

type HclConfig struct {
	DataDir        string `hcl:"data_dir,optional"`
	MetricsEnabled bool   `hcl:"metrics_enabled,optional"`
	LogLevel       string `hcl:"log_level,optional"`
	ListenAddress  string `hcl:"listen_address,optional"`
	// HttpInputEnabled bool
	Inputs          []HclConfigItem                  `hcl:"input,block"`
	Filters         []HclConfigItemWithSubscriptions `hcl:"filter,block"`
	Transformations []HclConfigItemWithSubscriptions `hcl:"transformation,block"`
	Destinations    []HclConfigItemWithSubscriptions `hcl:"destination,block"`
}

func processConfigItem(items []HclConfigItem, ctx *hcl.EvalContext) {
	for _, item := range items {
		ctx.Variables[item.Kind] = cty.ObjectVal(map[string]cty.Value{
			item.ID: cty.StringVal(item.ID),
		})
	}
}
func processConfigItemWithSubs(items []HclConfigItemWithSubscriptions, ctx *hcl.EvalContext) {
	for _, item := range items {
		ctx.Variables[item.Kind] = cty.ObjectVal(map[string]cty.Value{
			item.ID: cty.StringVal(item.ID),
		})
	}
}

func (hc HclConfig) EvalContext() *hcl.EvalContext {
	ctx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
	}
	// processConfigItem(hc.Inputs, ctx)
	// processConfigItemWithSubs(hc.Filters, ctx)
	// processConfigItemWithSubs(hc.Transformations, ctx)
	// processConfigItemWithSubs(hc.Destinations, ctx)
	return ctx
}

type HclConfigItem struct {
	Kind string   `hcl:"kind,label"`
	ID   string   `hcl:"id,label"`
	Body hcl.Body `hcl:",remain"`
}
type HclConfigItemWithSubscriptions struct {
	Kind          string   `hcl:"kind,label"`
	ID            string   `hcl:"id,label"`
	Subscriptions []string `hcl:"subscriptions"`
	Body          hcl.Body `hcl:",remain"`
}

func LoadHCL(fileName string, data []byte) (*Config, error) {
	var diags hcl.Diagnostics
	var config HclConfig
	if err := hclsimple.Decode(fileName, data, nil, &config); err != nil {
		slog.Error("error during simple decode", "error", err)
		return nil, err
	}
	opts := []ConfigOption{}
	if config.ListenAddress != "" {
		opts = append(opts, WithListenAddress(config.ListenAddress))
	}
	if config.LogLevel != "" {
		opts = append(opts, WithLogLevel(config.LogLevel))
	}
	if config.MetricsEnabled {
		opts = append(opts, WithMetricsEnabled())
	}
	ctx := config.EvalContext()
	for _, in := range config.Inputs {
		impl, more := input.HclImpl(in.Kind, ctx, in.Body)
		diags = append(diags, more...)
		if more.HasErrors() {
			continue
		}
		opts = append(opts, WithInput(in.ID, in.Kind, impl))
	}
	for _, fltr := range config.Filters {
		impl, more := filter.HclImpl(fltr.Kind, ctx, fltr.Body)
		diags = append(diags, more...)
		if more.HasErrors() {
			continue
		}
		opts = append(opts, WithFilter(fltr.ID, fltr.Kind, fltr.Subscriptions, impl))
	}
	for _, tr := range config.Transformations {
		impl, more := transformation.HclImpl(tr.Kind, tr.Body)
		diags = append(diags, more...)
		if more.HasErrors() {
			continue
		}
		opts = append(opts, WithTransformer(tr.ID, tr.Kind, tr.Subscriptions, impl))
	}
	for _, d := range config.Destinations {
		impl, more := destination.HclImpl(d.Kind, d.Body)
		diags = append(diags, more...)
		if more.HasErrors() {
			continue
		}
		opts = append(opts, WithDestination(d.ID, d.Kind, d.Subscriptions, impl))
	}

	if diags.HasErrors() {
		slog.Error("error in hcl diagnostics", "error", diags.Error())
		return nil, diags
	}
	cfg := New(opts...)
	return cfg, nil
}
func (c *Config) addSubscription(ref string, subscriber pubsub.Subscriber) {
	broker, ok := c.references[ref]
	if ok {
		broker.AddSubscriber(subscriber)
	}
}
func (c *Config) Init() {
	if c.ListenAddress == "" {
		c.ListenAddress = ":8080"
	}

	for _, d := range c.Destinations {
		for _, name := range d.Subscriptions {
			c.addSubscription(name, d.Subscriber)
		}
	}

	for _, t := range c.Transformations {
		for _, name := range t.Subscriptions {
			c.addSubscription(name, t.Subscriber)
		}
	}

	for _, f := range c.Filters {
		for _, name := range f.Subscriptions {
			c.addSubscription(name, f.Subscriber)
		}
	}

}
func New(opts ...ConfigOption) *Config {
	c := &Config{
		references:      make(map[string]pubsub.Broker),
		Inputs:          make([]*input.Input, 0),
		Filters:         make([]*filter.Filter, 0),
		Destinations:    make([]*destination.Destination, 0),
		Transformations: make([]*transformation.Transformation, 0),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Init()
	return c
}
