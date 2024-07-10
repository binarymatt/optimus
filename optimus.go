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
	cfg *config.Config
}

func New(cfg *config.Config) *Optimus {
	o := &Optimus{
		cfg: cfg,
	}
	return o
}

func (o *Optimus) Run(ctx context.Context) error {
	slog.Info("starting optimus runtime")
	cfg := o.cfg
	// read config
	eg, ctx := errgroup.WithContext(ctx)
	// input map
	parents := map[string]*pubsub.Broker{}
	slog.Debug("configuring inputs")
	for key, input := range cfg.Inputs {
		input.Init(key)
		parents[key] = input.Broker
	}

	slog.Debug("configuring filters")
	for key, filter := range cfg.Filters {
		filter.Init(key)
		parents[filter.ID] = filter.Broker
	}

	slog.Debug("configuring and starting destinations")
	for key, d := range cfg.Destinations {
		destination := d
		destination.ID = key
		// create destination channel
		if err := destination.Init(); err != nil {
			return err
		}
		// setup subscriptions
		for _, name := range destination.Subscriptions {
			broker := parents[name]
			slog.Debug("setting up destination parent", "name", name)
			broker.AddSubscriber(destination.Subscriber)
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
			broker := parents[name]
			broker.AddSubscriber(filter.Subscriber)
		}
		//start goroutine
		eg.Go(func() error {
			return filter.Process(ctx)
		})
	}

	slog.Debug("starting inputs")
	for _, i := range cfg.Inputs {
		input := i
		if err := input.Internal.Setup(ctx, input.Broker); err != nil {
			slog.Error("could not setup input", "input", input.ID, "kind", input.Kind)
			continue
		}
		eg.Go(func() error {
			return input.Process(ctx)
		})
	}

	if o.cfg.ListenAddress != "" {
		eg.Go(func() error {
			http.Handle("/metrics", promhttp.Handler())
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
				slog.Warn("done shutting down api server")
			}()

			slog.Info("api server starting")
			if err := server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					slog.Warn("server shutdown")
					return nil
				}
				slog.Error("got error during listen and serve", "error", err)
				return err
			}
			return nil
		})
	}
	slog.Info("everything starting, running...")
	return eg.Wait()
}
