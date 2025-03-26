package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Manager 管理配置文件的加载和监控
type Manager struct {
	viper  *viper.Viper
	mu     sync.RWMutex
	config interface{}
}

// NewConfigManager 创建一个新的 ConfigManager 实例
func NewConfigManager(configStruct interface{}) *Manager {
	return &Manager{
		viper:  viper.New(),
		config: configStruct,
	}
}

// LoadFile 加载配置文件到指定的结构体指针(结构体自由定义)
func (c *Manager) LoadFile(filePath string, watch bool) error {
	c.viper.SetConfigFile(filePath)

	if err := c.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.viper.Unmarshal(c.config); err != nil {
		return fmt.Errorf("failed to unmarshal config to struct: %w", err)
	}

	// 如果需要监控配置文件变化
	if watch {
		c.viper.WatchConfig()
		c.viper.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s", e.Name)
			c.mu.Lock()
			defer c.mu.Unlock()
			if err := c.viper.Unmarshal(c.config); err != nil {
				log.Printf("Failed to reload config: %v", err)
			}
		})
	}

	return nil
}

// GetConfig 返回当前配置
func (c *Manager) GetConfig() interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}
