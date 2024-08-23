package optimus

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	"github.com/binarymatt/optimus/config"
	"github.com/binarymatt/optimus/internal/pubsub"
)

type Optimus struct {
	cfg     *config.Config
	parents map[string]*pubsub.Broker
}

func New(cfg *config.Config) (*Optimus, error) {
	cfg.Init()
	o := &Optimus{
		cfg:     cfg,
		parents: make(map[string]*pubsub.Broker),
	}
	err := o.setup()
	return o, err
}
func (o *Optimus) setup() error {
	cfg := o.cfg
	// parents := map[string]*pubsub.Broker{}
	slog.Debug("configuring inputs")
	// input map
	for key, input := range cfg.Inputs {
		if input.Kind == "http" {
			o.cfg.HttpInputEnabled = true
		}
		broker, err := input.Init(key)
		if err != nil {
			return err
		}
		o.parents[key] = broker
	}
	slog.Debug("configuring filters")
	for key, filter := range cfg.Filters {
		broker := filter.Init(key)
		o.parents[key] = broker
	}
	for key, destination := range cfg.Destinations {
		if err := destination.Init(key); err != nil {
			slog.Error("could not initialize destiantion", "id", key)
			return err
		}
	}
	return nil
}
func (o *Optimus) Run(ctx context.Context) error {
	slog.Info("starting optimus runtime")
	cfg := o.cfg
	//if err := o.Setup(); err != nil {
	//	return err
	//}
	eg, ctx := errgroup.WithContext(ctx)

	slog.Debug("configuring and starting destinations")
	for _, d := range cfg.Destinations {
		destination := d
		// create destination channel
		// setup subscriptions
		for _, name := range destination.Subscriptions {
			broker, ok := o.parents[name]
			if ok {
				slog.Debug("setting up destination parent", "name", name)
				broker.AddSubscriber(destination.Subscriber)
			}
		}
		// start goroutine
		eg.Go(func() error {
			destination.Process(ctx)
			return nil
		})
	}
	slog.Debug("starting filters")
	for _, f := range cfg.Filters {
		filter := f
		// setup Subscriptions
		for _, name := range filter.Subscriptions {
			broker, ok := o.parents[name]
			if ok {
				broker.AddSubscriber(filter.Subscriber)
			}
		}
		slog.Debug("starting filter go routine", "name", filter.Kind)
		//start goroutine
		eg.Go(func() error {
			return filter.Process(ctx)
		})
	}

	slog.Debug("starting inputs")
	for _, i := range cfg.Inputs {
		input := i
		slog.Debug("starting input goroutine", "name", i.ID)
		eg.Go(func() error {
			return input.Process(ctx)
		})
	}

	if o.cfg.HttpInputEnabled || o.cfg.MetricsEnabled {
		eg.Go(func() error {
			slog.Debug("setting up http server")
			if o.cfg.MetricsEnabled {
				http.Handle("/metrics", promhttp.Handler())
			}
			server := http.Server{
				Addr:    o.cfg.ListenAddress,
				Handler: h2c.NewHandler(http.DefaultServeMux, &http2.Server{}),
			}

			go func() {
				<-ctx.Done()
				timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancelTimeout()

				if err := server.Shutdown(timeoutCtx); err != nil {
					slog.Error("error shutting service", "error", err)
				}
				slog.Debug("done shutting down api server")
			}()

			slog.Debug("api server starting")
			if err := server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					slog.Debug("server shutdown")
					return nil
				}
				slog.Error("got error during listen and serve", "error", err)
				return err
			}
			return nil
		})
	}
	slog.Info("all components started, running...")
	return eg.Wait()
}
