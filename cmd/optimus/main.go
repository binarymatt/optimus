package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"gopkg.in/yaml.v3"

	"github.com/binarymatt/optimus"
	"github.com/binarymatt/optimus/config"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "host",
			Value: "0.0.0.0",
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:  "port",
			Value: 3000,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "console",
			Required: false,
			Value:    false,
		}),
		&cli.StringFlag{
			Name: "config",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "log_level",
			Aliases: []string{"L"},
			Value:   "info",
		}),
	}

	app := &cli.App{
		Action: Run,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:  flags,
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("error running optimus", "error", err)
	}

}

func setupLogging(console bool, levelStr string) {
	var logger *slog.Logger

	var level slog.Leveler
	switch levelStr {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level, AddSource: true}
	if console {
		logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:     level,
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	slog.SetDefault(logger)
}
func loadConfig(filePath string) (*config.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	cfg := config.Config{}
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

func Run(cctx *cli.Context) error {
	setupLogging(cctx.Bool("console"), cctx.String("log_level"))
	slog.Info("starting optimus")
	// TODO - do i need signal context

	ctx, cancel := signal.NotifyContext(
		cctx.Context,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	// read config
	cfg, err := loadConfig(cctx.String("config"))
	if err != nil {
		return err
	}
	o := optimus.New(cfg)
	return o.Run(ctx)
}
