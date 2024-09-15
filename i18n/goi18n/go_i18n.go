package goi18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// I18nManager 全局并发安全
type I18nManager struct {
	bundle     *i18n.Bundle
	localizer  *i18n.Localizer
	localesDir string
	langCache  map[string]map[string]string
	mu         sync.RWMutex
}

// NewI18nManager 初始化并返回一个 I18nManager 实例
func NewI18nManager(localesDir string, defaultLangs ...string) (*I18nManager, error) {
	bundle := i18n.NewBundle(language.English) // 默认语言为英语
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	manager := &I18nManager{
		bundle:     bundle,
		localesDir: localesDir,
		langCache:  make(map[string]map[string]string),
	}

	// 加载 localesDir 中的所有 .yaml 文件
	files, err := filepath.Glob(filepath.Join(localesDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("无法读取语言文件: %w", err)
	}

	for _, file := range files {
		if _, err := bundle.LoadMessageFile(file); err != nil {
			return nil, fmt.Errorf("加载语言文件失败 %s: %w", file, err)
		}
	}

	if len(defaultLangs) == 0 {
		defaultLangs = []string{"en"} // 如果未指定默认语言，则使用英语
	}

	manager.localizer = i18n.NewLocalizer(bundle, defaultLangs...)

	return manager, nil
}

// GetTemplateLangMap 获取指定语言的模板映射
func (m *I18nManager) GetTemplateLangMap(lang string) (map[string]string, error) {
	m.mu.RLock()
	if langMap, exists := m.langCache[lang]; exists {
		m.mu.RUnlock()
		return langMap, nil
	}
	m.mu.RUnlock()

	langFilePath := filepath.Join(m.localesDir, lang+".yaml")
	data, err := os.ReadFile(langFilePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取语言文件: %w", err)
	}

	var langMap map[string]string
	if err := yaml.Unmarshal(data, &langMap); err != nil {
		return nil, fmt.Errorf("解析语言文件失败: %w", err)
	}

	m.mu.Lock()
	m.langCache[lang] = langMap
	m.mu.Unlock()

	return langMap, nil
}

// Translate 按照语言文件翻译数据中的字符串
func (m *I18nManager) Translate(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		translation, err := m.localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
		if err == nil && translation != "" {
			if _, ok := value.(string); ok {
				result[key] = translation
			} else {
				result[key] = value
			}
		} else {
			result[key] = value
		}
	}
	return result
}
