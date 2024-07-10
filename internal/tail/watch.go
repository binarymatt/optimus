package tail

import "github.com/fsnotify/fsnotify"

func Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	return nil
}
