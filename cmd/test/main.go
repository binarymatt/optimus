package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lmittmann/tint"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/binarymatt/optimus"
	"github.com/binarymatt/optimus/config"
	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

func loadConfig(filePath string) (*config.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return config.LoadHCL(filePath, data)
}
func main() {
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	slog.Warn("loading config")
	rawCfg, err := os.ReadFile("sample_config.hcl")
	if err != nil {
		return
	}
	cfg, err := config.LoadHCL("sample_config.hcl", rawCfg)
	if err != nil {
		return
	}
	c := make(chan *optimusv1.LogEvent)
	config.WithChannelInput("testing", c)(cfg)
	o := optimus.New(cfg)
	eg.Go(func() error {
		return o.Run(ctx)
	})
	eg.Go(func() error {
		for i := range 10 {
			data, err := structpb.NewStruct(map[string]interface{}{"test": "this", "id": i})
			if err != nil {
				return err
			}
			c <- &optimusv1.LogEvent{
				Id:   fmt.Sprintf("testing%d", i),
				Data: data,
			}
		}
		slog.Info("done adding records to channel")
		return nil
	})
	if err := eg.Wait(); err != nil {
		slog.Error("errorgroup wait error", "error", err)
	}
	slog.Warn("done with test command")
}

/*
func mainOld() {
	slog.Info("starting")
	path := "./tmp"

	//r := &recorder{
	// 	files: map[string]int64{},
	// }
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("error creating watcher", "error", err)
		os.Exit(1)
	}
	defer w.Close()
	slog.Info("watcher setup")
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case event, ok := <-w.Events:
				if !ok {
					slog.Error("channel closed")
					return nil
				}
				if event.Has(fsnotify.Create) {
					// r.files[event.Name] = 0
					slog.Info("file created")
				}
				if event.Has(fsnotify.Write) {
					slog.Debug("file written to", "event", event)
				}
				slog.Info("fs event", "event", event)
			}
		}
	})
	if err := w.Add(path); err != nil {
		slog.Error("error adding path", "error", err)
		return
	}
	slog.Info("running... ", "eg", eg)
	// _ = eg.Wait()
	f, err := os.Open("./tmp/second.json")
	if err != nil {
		slog.Error("cloudl not open")
		return
	}
	pos, err := f.Seek(0, io.SeekStart)
	slog.Error("done seeking", "pos", pos, "error", err)
	reader := bufio.NewReader(f)
	data, err := reader.ReadBytes('\n')
	slog.Error("read from pos", "data", data, "error", err)
	_, _ = f.Seek(int64(len(data)), io.SeekStart)
	reader.Reset(f)
	data, err = reader.ReadBytes('\n')
	slog.Error("read from pos", "data", data, "error", err)

}
*/
