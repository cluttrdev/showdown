package watch

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func WatchFile(filePath string, onWrite func()) (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating new watcher: %w", err)
	}

	// add file
	st, err := os.Lstat(filePath)
	if err != nil {
		return nil, fmt.Errorf("error requesting file info: %w", err)
	}

	if st.IsDir() {
		return nil, fmt.Errorf("%q is a directory, not a file", filePath)
	}

	err = w.Add(filePath)
	if err != nil {
		w.Close()
		return nil, fmt.Errorf("error adding file: %w", err)
	}

	// start listening for events
	go eventLoop(w, onWrite)

	return w, nil
}

func eventLoop(w *fsnotify.Watcher, onWrite func()) {
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

			if e.Has(fsnotify.Write) {
				onWrite()
			}
		}
	}
}
