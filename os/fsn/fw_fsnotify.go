package fsn

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// EventHandlers holds the functions to be called on various file system events.
type EventHandlers struct {
	OnCreate func(string)
	OnWrite  func(string)
	OnRemove func(string)
	OnRename func(string)
	OnChmod  func(string)
}

// Watcher monitors file system events using fsnotify.
type Watcher struct {
	watcher    *fsnotify.Watcher
	mu         sync.Mutex
	includeExt []string // File extensions to monitor

	// Event handlers
	OnCreate func(string)
	OnWrite  func(string)
	OnRemove func(string)
	OnRename func(string)
	OnChmod  func(string)
}

// NewFsnWatcher creates a new Watcher with the specified file extensions.
func NewFsnWatcher(includeExts []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw := &Watcher{
		watcher:    watcher,
		includeExt: includeExts,
	}
	return fw, nil
}

// StartWatcher initializes the watcher, adds directories, sets event handlers, and starts monitoring.
func StartWatcher(dirs []string, exts []string, handlers EventHandlers) (*Watcher, error) {
	watcher, err := NewFsnWatcher(exts)
	if err != nil {
		return nil, err
	}

	// Add directories recursively
	for _, dir := range dirs {
		err = watcher.AddDirRecursive(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to add directory recursively %s: %w", dir, err)
		}
		log.Printf("Started monitoring root directory: %s\n", dir)
	}

	// Register event handlers
	watcher.OnCreate = handlers.OnCreate
	watcher.OnWrite = handlers.OnWrite
	watcher.OnRemove = handlers.OnRemove
	watcher.OnRename = handlers.OnRename
	watcher.OnChmod = handlers.OnChmod

	// Start the watcher
	watcher.Start()

	return watcher, nil
}

// AddDirRecursive adds directories and subdirectories to the watch list recursively.
func (w *Watcher) AddDirRecursive(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.watchFileLocked(subPath)
		}
		return nil
	})
}

// watchFileLocked watches a single directory or file. Assumes the mutex is already locked.
func (w *Watcher) watchFileLocked(path string) error {
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create path %s: %w", path, err)
		}
		log.Printf("Path %s did not exist, created it\n", path)
	} else if err != nil {
		return fmt.Errorf("unable to access path %s: %w", path, err)
	}

	// Add the path to the watcher
	err := w.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("unable to watch path %s: %w", path, err)
	}
	log.Printf("Watching: %s\n", path)
	return nil
}

// shouldMonitor checks if a file should be monitored based on its extension.
func (w *Watcher) shouldMonitor(filename string) bool {
	// Check if it's a directory
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}

	// Check if the file has one of the included extensions
	ext := strings.ToLower(filepath.Ext(filename))
	for _, includeExt := range w.includeExt {
		if ext == includeExt {
			return true
		}
	}
	return false
}

// Start begins the main loop for file monitoring and event handling.
func (w *Watcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// Before processing, check if the file should be monitored
				if !w.shouldMonitor(event.Name) {
					continue
				}

				// Handle the event
				w.handleEvent(event)

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v\n", err)
			}
		}
	}()
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Handle create events
	if event.Op&fsnotify.Create == fsnotify.Create {
		fi, err := os.Stat(event.Name)
		if err == nil {
			if fi.IsDir() {
				// Recursively add new directories to the watch list
				if err := w.AddDirRecursive(event.Name); err != nil {
					log.Printf("Failed to recursively watch new directory %s: %v\n", event.Name, err)
				}
				// Walk the directory and trigger OnCreate for files
				filepath.Walk(event.Name, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						log.Printf("Failed to access path %s: %v\n", path, err)
						return nil
					}
					if !info.IsDir() && w.shouldMonitor(path) {
						if w.OnCreate != nil {
							go func(p string) { w.OnCreate(p) }(path)
						}
					}
					return nil
				})
			} else if w.shouldMonitor(event.Name) {
				// Handle creation of new files
				if w.OnCreate != nil {
					go func() { w.OnCreate(event.Name) }()
				}
			}
		} else {
			log.Printf("Failed to get file info %s: %v\n", event.Name, err)
		}
		return
	}

	// Handle other events
	if w.shouldMonitor(event.Name) {
		if event.Op&fsnotify.Write == fsnotify.Write && w.OnWrite != nil {
			go func() { w.OnWrite(event.Name) }()
		}

		if event.Op&fsnotify.Remove == fsnotify.Remove && w.OnRemove != nil {
			go func() { w.OnRemove(event.Name) }()
		}

		if event.Op&fsnotify.Rename == fsnotify.Rename && w.OnRename != nil {
			go func() { w.OnRename(event.Name) }()
		}

		if event.Op&fsnotify.Chmod == fsnotify.Chmod && w.OnChmod != nil {
			go func() { w.OnChmod(event.Name) }()
		}
	}
}

// Close stops the file watcher and releases resources.
func (w *Watcher) Close() error {
	return w.watcher.Close()
}
