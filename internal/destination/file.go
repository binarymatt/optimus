package destination

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

type FileDestination struct {
	Path string `yaml:"path"`
	w    io.Writer
	f    *os.File
}

var (
	ErrMissingPath = errors.New("missing path for file destination")
)

func (f *FileDestination) Setup(cfg map[string]any) error {
	slog.Debug("file setup", "cfg", cfg)
	val, ok := cfg["path"]
	if !ok {
		slog.Error("could not setup file destination, missing path")
		return ErrMissingPath
	}
	f.Path = val.(string)

	openFile, err := os.OpenFile(f.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("could not open file", "error", err)
		return err
	}
	f.f = openFile
	f.w = bufio.NewWriter(openFile)
	return nil
}

func (f *FileDestination) Deliver(ctx context.Context, event *optimusv1.LogEvent) error {
	raw, err := json.Marshal(event.Data.AsMap())
	if err != nil {
		slog.Error("could not marshall data for writing to file", "error", err)
		return err
	}
	_, err = f.f.WriteString(string(raw) + "\n")
	return err
}
