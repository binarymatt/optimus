package destination

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type StdOutDestination struct {
	encoder *json.Encoder
}

func (sd *StdOutDestination) Setup() error {
	sd.encoder = json.NewEncoder(os.Stdout)
	return nil
}

func (sd *StdOutDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	if err := sd.encoder.Encode(event.Data.AsMap()); err != nil {
		slog.Error("could not encode data to stdout", "error", err)
		return err
	}
	return nil
}
