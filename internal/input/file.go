package input

import (
	"context"
	"io"
	"os"

	"github.com/binarymatt/optimus/internal/pubsub"
)

type FileInput struct {
	Path string `yaml:"path"`
}

func (fi *FileInput) Setup(ctx context.Context, broker *pubsub.Broker) error {
	return nil
}

func (fi *FileInput) Process(ctx context.Context) error {
	return nil
}
