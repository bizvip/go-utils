package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Manager 通用配置管理器，使用Go 1.24泛型提供类型安全
type Manager[T any] struct {
	viper  *viper.Viper
	mu     sync.RWMutex
	config *T
}

// Config 配置接口，允许自定义默认值
type Config interface {
	SetDefaults()
}

// New 创建新的泛型配置管理器
func New[T any](config *T) *Manager[T] {
	return &Manager[T]{
		viper:  viper.New(),
		config: config,
	}
}

// NewFromExecutable 自动从可执行文件同级目录加载配置
func NewFromExecutable[T any](config *T, configName string) (*Manager[T], error) {
	manager := New(config)

	configPath, err := getExecutableConfigPath(configName)
	if err != nil {
		return nil, fmt.Errorf("获取配置文件路径失败: %w", err)
	}

	return manager, manager.LoadFile(configPath, false)
}

// NewFromExecutableWithWatch 自动从可执行文件同级目录加载配置并监听变化
func NewFromExecutableWithWatch[T any](config *T, configName string) (*Manager[T], error) {
	manager := New(config)

	configPath, err := getExecutableConfigPath(configName)
	if err != nil {
		return nil, fmt.Errorf("获取配置文件路径失败: %w", err)
	}

	return manager, manager.LoadFile(configPath, true)
}

// LoadFile 加载配置文件到指定的结构体指针
func (m *Manager[T]) LoadFile(filePath string, watch bool) error {
	m.viper.SetConfigFile(filePath)

	if err := m.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 设置默认值（如果配置实现了Config接口）
	if configWithDefaults, ok := any(m.config).(Config); ok {
		configWithDefaults.SetDefaults()
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.viper.Unmarshal(m.config); err != nil {
		return fmt.Errorf("配置反序列化失败: %w", err)
	}

	if watch {
		m.viper.WatchConfig()
		m.viper.OnConfigChange(func(e fsnotify.Event) {
			log.Info().Str("file", e.Name).Msg("配置文件已更改")
			m.mu.Lock()
			defer m.mu.Unlock()
			if err := m.viper.Unmarshal(m.config); err != nil {
				log.Error().Err(err).Msg("重新加载配置失败")
			}
		})
	}

	return nil
}

// GetConfig 返回类型安全的配置
func (m *Manager[T]) GetConfig() *T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// UpdateConfig 原子性更新配置
func (m *Manager[T]) UpdateConfig(updateFn func(*T)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	updateFn(m.config)
}

// getExecutableConfigPath 获取可执行文件同级目录的配置文件路径
func getExecutableConfigPath(configName string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, configName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", fmt.Errorf("配置文件不存在: %s", configPath)
	}

	return configPath, nil
}
