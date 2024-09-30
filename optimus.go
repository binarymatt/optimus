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
)

type Optimus struct {
	cfg *config.Config
	// parents map[string]pubsub.Broker
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
	eg, ctx := errgroup.WithContext(ctx)

	slog.Debug("configuring and starting destinations")
	for _, d := range cfg.Destinations {
		destination := d
		// start goroutine
		eg.Go(func() error {
			destination.Process(ctx)
			return nil
		})
	}
	slog.Debug("starting transformations")
	for _, t := range cfg.Transformations {
		transformation := t
		slog.Debug("starting transformation go routine", "name", transformation.ID)
		eg.Go(func() error {
			return transformation.Process(ctx)
		})
	}
	slog.Debug("starting filters")
	for _, f := range cfg.Filters {
		filter := f
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
	// TODO add transformations
	slog.Info("all components started, running...")
	return eg.Wait()
}
