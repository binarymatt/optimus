package input

import (
	"context"
	"encoding/json"
	"log/slog"
	"path/filepath"

	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/structpb"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
	"github.com/binarymatt/optimus/internal/metrics"
	"github.com/binarymatt/optimus/internal/pubsub"
	"github.com/binarymatt/optimus/internal/tail"
)

type FileInput struct {
	ID      string
	Path    string `yaml:"path"`
	tracker *tail.TailTracker
}

func (fi *FileInput) Initialize(id string, broker *pubsub.Broker) error {
	pathName := filepath.Clean(fi.Path)
	slog.Debug("setting up file input", "path", pathName)
	t, err := tail.NewTracker()
	if err != nil {
		return err
	}
	fi.tracker = t
	_, err = fi.tracker.AddPath(fi.Path, func(path, line string) error {
		uid := ulid.Make()
		var log map[string]interface{}
		if err := json.Unmarshal([]byte(line), &log); err != nil {
			slog.Error("could not unmarshal line", "error", err)
			return err
		}
		data, err := structpb.NewStruct(log)
		if err != nil {
			return err
		}

		event := &optimusv1.LogEvent{
			Id:     uid.String(),
			Source: path,
			Data:   data,
		}
		metrics.IncProcessedRecord("file_input", id)
		broker.Broadcast(event)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (fi *FileInput) Process(ctx context.Context) error {
	return fi.tracker.Run(ctx)
}
