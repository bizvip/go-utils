/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package config

import (
	"fmt"

	goutils "github.com/bizvip/go-utils"
	"github.com/spf13/viper"
)

// LoadConfig 加载指定配置文件并解析到指定的结构体指针中
func LoadConfig(configPath string, configStruct interface{}) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("failed to unmarshal config to struct: %w", err)
	}

	return nil
}

// WriteKeyValue 将数据写入指定配置文件中的指定层级
func WriteKeyValue(configPath string, key string, value interface{}) error {
	viper.SetConfigFile(configPath)

	// 读取现有配置
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 设置指定key的值
	viper.Set(key, value)

	// 将Viper配置写入文件
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// WriteConfigFile 将数据写入指定配置文件
func WriteConfigFile(configPath string, configStruct interface{}) error {
	viper.SetConfigFile(configPath)

	// 将结构体内容写入 Viper
	sm, err := goutils.NewStructUtils().StructToMap(configStruct)
	if err != nil {
		return err
	}
	if err = viper.MergeConfigMap(sm); err != nil {
		return fmt.Errorf("failed to merge config map: %w", err)
	}

	// 将 Viper 配置写入文件
	if err = viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

//
//// WriteConfigBigCache 将数据写入 BigCache 和配置文件
//func WriteConfigBigCache(configPath string, cache *bigcache.BigCache, configStruct interface{}) error {
//	// 将结构体转换为 map
//	sm, err := goutils.NewStructUtils().StructToMap(configStruct)
//	if err != nil {
//		return fmt.Errorf("failed to convert struct to map: %w", err)
//	}
//
//	// 将 map 转换为 JSON
//	jsonData, err := json.Marshal(sm)
//	if err != nil {
//		return fmt.Errorf("failed to marshal config to json: %w", err)
//	}
//
//	// 将 JSON 数据写入 BigCache
//	if err := cache.Set("config", jsonData); err != nil {
//		return fmt.Errorf("failed to write config to bigcache: %w", err)
//	}
//
//	// 将配置写入配置文件
//	viper.SetConfigFile(configPath)
//	if err := viper.MergeConfigMap(sm); err != nil {
//		return fmt.Errorf("failed to merge config map: %w", err)
//	}
//	if err := viper.WriteConfigAs(configPath); err != nil {
//		return fmt.Errorf("failed to write config file: %w", err)
//	}
//
//	return nil
//}
