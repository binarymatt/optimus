package destination

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type StdOutDestination struct {
	encoder *json.Encoder
	Writer  io.Writer
}

func (sd *StdOutDestination) Setup() error {
	var writer io.Writer
	if sd.Writer != nil {
		writer = sd.Writer
	} else {
		writer = os.Stdout
	}
	sd.encoder = json.NewEncoder(writer)
	return nil
}

func (sd *StdOutDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	if err := sd.encoder.Encode(event.Data.AsMap()); err != nil {
		slog.Error("could not encode data to stdout", "error", err)
		return err
	}
	return nil
}
func (sd *StdOutDestination) Close() error {
	return nil
}
