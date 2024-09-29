package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type FileDestination struct {
	Path string `hcl:"path"`
	w    io.Writer
	f    *os.File
}

var (
	ErrMissingPath = errors.New("missing path for file destination")
)

func (f *FileDestination) Setup() error {
	slog.Debug("file setup", "path", f.Path)
	if f.Path == "" {

		return ErrMissingPath
	}

	openFile, err := os.OpenFile(f.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("could not open file", "error", err)
		return err
	}
	f.f = openFile
	f.w = openFile
	return nil
}

func (f *FileDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	raw, err := json.Marshal(event.Data.AsMap())
	if err != nil {
		slog.Error("could not marshall data for writing to file", "error", err)
		return err
	}
	combined := bytes.NewBuffer(raw)
	combined.WriteString("\n")
	_, err = f.w.Write(combined.Bytes())

	return err
}
func (f *FileDestination) Close() error {
	return f.f.Close()
}
