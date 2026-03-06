package configer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Manager[T any] struct {
	path    string
	options []LoadOption[T]

	mu      sync.RWMutex
	current *T
}

func NewManager[T any](path string, options ...LoadOption[T]) *Manager[T] {
	return &Manager[T]{
		path:    strings.TrimSpace(path),
		options: options,
	}
}

func (m *Manager[T]) Load() (*T, error) {
	cfg, err := Load[T](m.path, m.options...)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	m.current = cfg
	m.mu.Unlock()
	return cfg, nil
}

func (m *Manager[T]) Reload() (*T, error) {
	return m.Load()
}

func (m *Manager[T]) Current() *T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

func (m *Manager[T]) Watch(ctx context.Context, onChange func(*T), onError func(error)) error {
	if strings.TrimSpace(m.path) == "" {
		return fmt.Errorf("manager path is required")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher failed: %w", err)
	}

	dir := filepath.Dir(m.path)
	base := filepath.Base(m.path)
	if err := watcher.Add(dir); err != nil {
		_ = watcher.Close()
		return fmt.Errorf("watch config directory %q failed: %w", dir, err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) != base {
					continue
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
					continue
				}
				cfg, err := m.Reload()
				if err != nil {
					if onError != nil {
						onError(err)
					}
					continue
				}
				if onChange != nil {
					onChange(cfg)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if onError != nil {
					onError(err)
				}
			}
		}
	}()

	return nil
}
