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

// Watcher 根据系统文件变动事件统治来监控文件
type Watcher struct {
	watcher    *fsnotify.Watcher
	mu         sync.Mutex
	includeExt []string // 只监控的文件后缀名

	// 事件处理函数
	OnCreate func(string)
	OnWrite  func(string)
	OnRemove func(string)
	OnRename func(string)
	OnChmod  func(string)
}

// NewFsnWatcher 添加要监控的目录，支持递归子目录
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

// AddDirRecursive 递归地添加目录及其子目录到监控列表
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

// shouldMonitor 检查文件是否应当被监控
func (w *Watcher) shouldMonitor(filename string) bool {
	// 检查是否是目录
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}

	// 检查文件是否有指定的后缀名
	ext := strings.ToLower(filepath.Ext(filename))
	for _, includeExt := range w.includeExt {
		if ext == includeExt {
			return true
		}
	}
	return false
}

// WatchFile 开始监控文件或目录 如果监控的目录不存在，则会自动创建该目录
func (w *Watcher) WatchFile(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.watchFileLocked(path)
}

// 内部方法，假定已获得锁
func (w *Watcher) watchFileLocked(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 如果路径不存在，尝试创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("无法创建路径 %s: %w", path, err)
		}
		log.Printf("路径 %s 不存在，已创建\n", path)
	} else if err != nil {
		return fmt.Errorf("无法访问路径 %s: %w", path, err)
	}

	// 添加监控路径
	err := w.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("无法监控路径 %s: %w", path, err)
	}
	log.Printf("正在监控: %s\n", path)
	return nil
}

// Start 启动文件监控的主循环，并处理监控事件
func (w *Watcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// 在处理所有事件前，检查文件是否应当被监控
				if !w.shouldMonitor(event.Name) {
					log.Printf("忽略文件: %s\n", event.Name)
					continue
				}

				// 处理事件
				w.handleEvent(event)

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("监控错误: %v\n", err)
			}
		}
	}()
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	// 如果是创建事件
	if event.Op&fsnotify.Create == fsnotify.Create {
		fi, err := os.Stat(event.Name)
		if err == nil {
			if fi.IsDir() {
				// 递归添加新目录到监控列表
				if err := w.AddDirRecursive(event.Name); err != nil {
					log.Printf("无法递归监控新目录 %s: %v\n", event.Name, err)
				}
			} else if w.shouldMonitor(event.Name) {
				// 处理新文件的创建事件
				if w.OnCreate != nil {
					go func() { w.OnCreate(event.Name) }()
				}
			}
		} else {
			log.Printf("无法获取文件信息 %s: %v\n", event.Name, err)
		}
		return
	}

	// 其他事件处理逻辑
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

// Close 停止文件监控并释放资源
func (w *Watcher) Close() error {
	return w.watcher.Close()
}
