package filewatcher

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// FileWatcher 是一个封装文件监控操作的类
type FileWatcher struct {
	watcher *fsnotify.Watcher
	mu      sync.Mutex

	// 事件处理函数
	OnCreate func(string)
	OnWrite  func(string)
	OnRemove func(string)
	OnRename func(string)
	OnChmod  func(string)
}

// NewFileWatcher 创建一个新的文件监控器
func NewFileWatcher() (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &FileWatcher{
		watcher: watcher,
	}, nil
}

// WatchFile 开始监控文件或目录 如果监控的目录不存在，则会自动创建该目录
func (fw *FileWatcher) WatchFile(path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

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
	err := fw.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("无法监控路径 %s: %w", path, err)
	}
	log.Printf("正在监控: %s\n", path)
	return nil
}

// Start 启动文件监控的主循环，并处理监控事件
func (fw *FileWatcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return
				}
				fw.handleEvent(event) // 处理事件

			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("监控错误: %v\n", err)
			}
		}
	}()
}

// handleEvent 处理监控到的文件系统事件，事件函数全部在 goroutine 中执行
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create && fw.OnCreate != nil {
		go fw.OnCreate(event.Name)
	}
	if event.Op&fsnotify.Write == fsnotify.Write && fw.OnWrite != nil {
		go fw.OnWrite(event.Name)
	}
	if event.Op&fsnotify.Remove == fsnotify.Remove && fw.OnRemove != nil {
		go fw.OnRemove(event.Name)
	}
	if event.Op&fsnotify.Rename == fsnotify.Rename && fw.OnRename != nil {
		go fw.OnRename(event.Name)
	}
	if event.Op&fsnotify.Chmod == fsnotify.Chmod && fw.OnChmod != nil {
		go fw.OnChmod(event.Name)
	}
}

// Close 停止文件监控并释放资源
func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}
