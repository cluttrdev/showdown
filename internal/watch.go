package watch

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

func WatchFile(filePath string, handler func(e fsnotify.Event)) (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Errorf("creating a new watcher: %v", err)
	}

	// start listening for events
	go eventLoop(w, handler)

	// add file
	st, err := os.Lstat(filePath)
	if err != nil {
		return nil, errors.Errorf("requesting file info: %v", err)
	}

	if st.IsDir() {
		return nil, errors.Errorf("%q is a directory, not a file", filePath)
	}

	err = w.Add(filePath)
	if err != nil {
		w.Close()
		return nil, errors.Errorf("adding file: %v", err)
	}

	return w, nil
}

func eventLoop(w *fsnotify.Watcher, handler func(e fsnotify.Event)) {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok { // Channel is closed (i.e. Watcher.Close() was called)
				return
			}
			log.Fatal(err)
		case e, ok := <-w.Events:
			if !ok { // Channel is closed (i.e. Watcher.Close() was called)
				return
			}

			handler(e)
		}
	}
}
