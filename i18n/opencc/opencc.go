package opencc

import (
	"fmt"
	"sync"

	"github.com/longbridgeapp/opencc"
)

// 常用模式别名
const (
	SchemaT2S  = "t2s"  // 繁体 -> 简体（大陆）
	SchemaS2T  = "s2t"  // 简体 -> 繁体（通用）
	SchemaS2TW = "s2tw" // 简体 -> 繁体（台湾正体）
	SchemaTW2S = "tw2s" // 台湾正体 -> 简体
	SchemaT2HK = "t2hk" // 繁体 -> 繁体（香港）
	SchemaHK2S = "hk2s" // 香港繁体 -> 简体
)

// 全局缓存：模式 -> 转换器
var (
	once       sync.Once
	initErr    error
	mu         sync.RWMutex
	converters map[string]*opencc.OpenCC
)

// 进行默认初始化：预热 t2s/s2t，其他模式按需创建。
// 注意：只在第一次使用时触发；如需显式预热，可调用 WarmUp。
func ensureInit() error {
	once.Do(func() {
		converters = make(map[string]*opencc.OpenCC, 4)

		// 预热常用模式（失败不 panic，而是记录错误）
		if c, err := opencc.New(SchemaT2S); err != nil {
			initErr = fmt.Errorf("opencc init %s failed: %w", SchemaT2S, err)
			return
		} else {
			converters[SchemaT2S] = c
		}

		if c, err := opencc.New(SchemaS2T); err != nil {
			initErr = fmt.Errorf("opencc init %s failed: %w", SchemaS2T, err)
			return
		} else {
			converters[SchemaS2T] = c
		}
	})
	return initErr
}

// WarmUp 允许在应用启动阶段显式预热若干模式（可选）。
// 例如：opencc.WarmUp(opencc.SchemaS2TW, opencc.SchemaT2HK)
func WarmUp(schemas ...string) error {
	if err := ensureInit(); err != nil {
		return err
	}
	for _, sc := range schemas {
		if _, err := getConverter(sc); err != nil {
			return err
		}
	}
	return nil
}

// 获取或创建指定模式的转换器（并发安全，带缓存）。
func getConverter(schema string) (*opencc.OpenCC, error) {
	if err := ensureInit(); err != nil {
		return nil, err
	}

	// 读锁优先
	mu.RLock()
	c := converters[schema]
	mu.RUnlock()
	if c != nil {
		return c, nil
	}

	// 不存在则写锁创建（双检）
	mu.Lock()
	defer mu.Unlock()

	if c = converters[schema]; c != nil {
		return c, nil
	}
	nc, err := opencc.New(schema)
	if err != nil {
		return nil, fmt.Errorf("opencc new(%s) failed: %w", schema, err)
	}
	converters[schema] = nc
	return nc, nil
}

// Convert 使用指定 OpenCC 模式进行转换。
// schema 示例：t2s/s2t/s2tw/tw2s/t2hk/hk2s 等。
func Convert(text, schema string) (string, error) {
	c, err := getConverter(schema)
	if err != nil {
		return "", err
	}
	return c.Convert(text)
}

// --------- 常用模式的便捷函数 ---------

// TradToSimp 繁体 -> 简体（大陆规范）。
func TradToSimp(text string) (string, error) {
	return Convert(text, SchemaT2S)
}

// SimpToTrad 简体 -> 繁体（通用）。
func SimpToTrad(text string) (string, error) {
	return Convert(text, SchemaS2T)
}

// SimpToTW 简体 -> 繁体（台湾正体）。
func SimpToTW(text string) (string, error) {
	return Convert(text, SchemaS2TW)
}

// TWToS 台湾正体 -> 简体。
func TWToS(text string) (string, error) {
	return Convert(text, SchemaTW2S)
}
