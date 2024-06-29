/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var mu sync.RWMutex

// LoadFile 加载配置文件到指定的结构体指针(结构体自由定义)
func LoadFile(filePath string, configStruct interface{}, watch bool) error {
	v := viper.New()
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if err := v.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("failed to unmarshal config to struct: %w", err)
	}

	// 如果需要监控配置文件变化
	if watch {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s", e.Name)
			mu.Lock()
			if err := v.Unmarshal(configStruct); err != nil {
				log.Printf("Failed to reload config: %v", err)
			}
			mu.Unlock()
		})
	}

	return nil
}
