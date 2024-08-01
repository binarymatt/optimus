package tail

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	ErrChanClosedEmpty = errors.New("channel is closed or empty")
)

// type Config struct{}
type LineProcessor = func(string, string) error

func getPathPos(path string) (int64, error) {
	slog.Debug("geting pos from file", "path", path)
	pos := int64(0)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		slog.Error("could not read file to get pos", "path", path)
		return pos, err
	}
	s := strings.TrimSuffix(string(data), "\n")
	return strconv.ParseInt(s, 10, 0)
}
func savePos(path string, pos int64) error {
	raw := strconv.FormatInt(pos, 10)
	return os.WriteFile(path, []byte(raw), 0644)
}
func checkpointFilePath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write([]byte(path))
	raw := h.Sum(nil)

	name := filepath.Join(homeDir, ".config", "optimus", fmt.Sprintf("%x.checkpoint", raw))
	if err := os.MkdirAll(filepath.Dir(name), os.ModePerm); err != nil {
		slog.Error("could not make config dir")
		return "", err
	}
	slog.Debug("returning checkpoint", "path", name)
	return name, nil
}
func NewInfo(path string, processor LineProcessor) (*PathInfo, error) {
	name := filepath.Clean(path)
	isD := isDir(name)
	chkPath, err := checkpointFilePath(name)
	if err != nil {
		slog.Error("checkpoint path resulted in error, not setting", "error", err)
	}

	info := PathInfo{
		path:      name,
		isDir:     isD,
		processor: processor,
	}

	slog.Debug("setting up info struct", "path", name, "is_dir", isD, "checkpoint", chkPath)
	if !isD {
		info.checkpoint = chkPath
		info.Events = make(chan fsnotify.Event)
		pos, err := getPathPos(chkPath)
		if err != nil {
			slog.Error("could not get pos", "error", err, "is", errors.Is(err, os.ErrNotExist), "path", chkPath)
			return nil, err
		}
		info.pos = pos
		f, err := os.Open(name)
		info.done = make(chan bool)
		if err != nil {
			return nil, err
		}
		_, err = f.Seek(info.pos, io.SeekStart)
		if err != nil {
			f.Close()
			return nil, err
		}
		info.file = f
		info.reader = bufio.NewReader(f)
	}
	return &info, nil

}

type PathInfo struct {
	isDir  bool
	path   string
	Events chan fsnotify.Event
	// tail      *tailer
	processor  LineProcessor
	pos        int64
	reader     *bufio.Reader
	file       *os.File
	done       chan bool
	checkpoint string
	// parentPath string
}

func (p *PathInfo) Done() {
	p.file.Close()
	close(p.done)
	close(p.Events)
}
func (p *PathInfo) Run() {
	p.ReadUntilEOF()
	for {
		select {
		case <-p.done:
			p.Done()
			return
		case evt, ok := <-p.Events:
			if !ok {
				return
			}
			if evt.Has(fsnotify.Write) {
				p.ReadUntilEOF()
				return
			}
		}
	}
}

// TODO - batch reads
func (p *PathInfo) ReadLines() ([]string, error) {
	var lines []string
	pos := p.pos
	slog.Debug("starting to read at position", "pos", pos)
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		pos = pos + int64(len([]byte(line)))
		lines = append(lines, line)
	}
	_, _ = p.file.Seek(pos, io.SeekStart)
	if err := savePos(p.checkpoint, pos); err != nil {
		slog.Error("could not save checkpoint", "error", err, "chk", p.checkpoint, "pos", pos)
		return nil, err
	}
	p.pos = pos
	p.reader.Reset(p.file)
	slog.Debug("read lines", "count", len(lines), "file", p.path)
	return lines, nil
}

func (p *PathInfo) ReadUntilEOF() {
	slog.Debug("reading from file", "name", p.path)
	lines, err := p.ReadLines()
	if err != nil {
		return
	}
	for _, line := range lines {
		if err := p.processor(p.path, line); err != nil {
			return
		}
	}
}

func NewTracker() (*TailTracker, error) {
	// eg, ctx := errgroup.WithContext(ctx)
	// ctx, cancel := context.WithCancel(ctx)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		// cancel()
		return nil, err
	}
	return &TailTracker{
		mux:   sync.Mutex{},
		paths: make(map[string]*PathInfo),
		//eg:      eg,
		// ctx:        ctx,
		watcher: w,
		// cancelFunc: cancel,
	}, nil
}

type TailTracker struct {
	// files map[string]
	mux     sync.Mutex
	paths   map[string]*PathInfo
	watcher *fsnotify.Watcher
	// eg         *errgroup.Group
	// ctx        context.Context
	// cancelFunc context.CancelFunc
}

func (t *TailTracker) Run(ctx context.Context) error {
	defer t.watcher.Close()
	// defer t.cancelFunc()
	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-t.watcher.Events:
			if !ok {
				return ErrChanClosedEmpty
			}
			t.processEvent(event)
		}
	}
}
func (t *TailTracker) AddFromParent(parent, name string) (*PathInfo, error) {
	t.mux.Lock()
	parentInfo, ok := t.paths[parent]
	t.mux.Unlock()
	if !ok {
		slog.Error("skipping creation, no parent info", "path", name)
		return nil, nil
	}
	slog.Debug("setting up with parent processor", "info", parentInfo)

	return t.AddPath(name, parentInfo.processor)
}

func (t *TailTracker) processEvent(event fsnotify.Event) {
	path := event.Name
	name := filepath.Clean(path)
	parent := filepath.Dir(name)

	slog.Debug("event processed", "operation", event.Op, "name", name, "path", path, "parent", parent)
	if event.Has(fsnotify.Create) {
		_, err := t.AddFromParent(parent, name)
		if err != nil {
			return
		}
	}
	if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
		t.mux.Lock()
		info, ok := t.paths[name]
		if ok {
			info.Done()
		}
		delete(t.paths, name)
		t.mux.Unlock()
	}
	if event.Has(fsnotify.Write) {
		var err error
		slog.Info("sending event to tailer")
		t.mux.Lock()
		info, ok := t.paths[name]
		t.mux.Unlock()
		if !ok {
			info, err = t.AddFromParent(parent, name)
			if err != nil {
				return
			}
		}
		info.Events <- event
		slog.Debug("what is info", "info", info)
	}
}

// AddPath adds the given path to the set of watched paths. If the path is a directory,
// files included in that directory will be added to the watcher as they are updated/created.
func (t *TailTracker) AddPath(path string, processor LineProcessor) (*PathInfo, error) {
	slog.Debug("adding path to tracker", "path", path)

	name := filepath.Clean(path)
	info, err := NewInfo(path, processor)
	if err != nil {
		return nil, err
	}
	slog.Debug("adding info struct to tracker", "path", name)
	t.mux.Lock()

	t.paths[name] = info
	if err := t.watcher.Add(path); err != nil {
		slog.Error("could not add to watcher", "error", err)
		t.mux.Unlock()
		return nil, err
	}
	t.mux.Unlock()
	if !info.isDir {
		slog.Debug("setting up goroutine to follow file")
		go info.Run()
	} else {
		slog.Debug("adding subpaths")
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				_, _ = t.AddFromParent(name, filepath.Join(path, entry.Name()))
			}
		}
	}
	return info, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
