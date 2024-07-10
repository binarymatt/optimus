package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
)
var (
	global *tracker
	goRun = func (){
		global = &tracker{
			mux: sync.Mutex{},
			paths: make(map[string]pathInfo)
		}
	}
)

type pathInfo struct {
	isDir  bool
	path   string
	events chan fsnotify.Event
}

type tracker struct {
	// files map[string]
	mux sync.Mutex
	paths map[string]pathInfo
}




type SeekInfo struct {
	Offset int64
	Whence int
}
type tailer struct {
	pos      int64
	filename string
	file     *os.File
	reader   *bufio.Reader
	// watcher  *fsnotify.Watcher
	Events chan fsnotify.Event
	lock   sync.Mutex
}

func NewTail(filename string) (*tailer, error) {
}
func (t *tailer) tailFile() {}
func (t *tailer) readline() (string, error) {
	t.lock.Lock()
	data, err := t.reader.ReadString('\n')
	return data, err
}
func (t *tailer) updatePosition(additionalPosData int64) error {
	t.pos = t.pos + additionalPosData
	_, err := t.file.Seek(t.pos, io.SeekStart)
	if err != nil {
		return err
	}
	t.reader.Reset(t.file)
	return nil
}
func (t *tailer) wait() error {
	select {
	case event, ok := <-t.Events:
		if !ok {
			return ErrChannelClosed
		}

	}
}
func (t *tailer) process() {
	for {
		line, err := t.readline()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// wait for changes
				continue
			}
			return
		}
	}
}

func NewTail(path string) (*tailer, error) {
	f, err := os.Open(path)
	if err != nil {
		slog.Error("could not open file for read", "error", err)
		return nil, err
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("could not create watcher", "error", err)
		return nil, err
	}
	w.Add(path)
	return &tailer{
		filename: path,
		file:     f,
		reader:   bufio.NewReader(f),
		watcher:  w,
	}, nil
}

func main() {
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

func openFile(path string) {
	// return os.Open(path)
}
