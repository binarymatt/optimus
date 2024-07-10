package destination

import (
	"context"
	"encoding/json"
	"os"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type StdOutDestination struct {
	encoder *json.Encoder
}

func (sd *StdOutDestination) Setup(cfg map[string]any) error {
	sd.encoder = json.NewEncoder(os.Stdout)
	return nil
}

func (sd *StdOutDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	if err := sd.encoder.Encode(event.Data.AsMap()); err != nil {
		return err
	}
	return nil
}
