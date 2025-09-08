package conf_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bizvip/go-utils/conf"
)

// TestAppConfig 测试配置结构
type TestAppConfig struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Debug   bool   `toml:"debug"`
}

func (c *TestAppConfig) SetDefaults() {
	if c.Name == "" {
		c.Name = "TestApp"
	}
	if c.Version == "" {
		c.Version = "1.0.0"
	}
}

func TestConfigManager_Basic(t *testing.T) {
	cfg := &TestAppConfig{}
	manager := conf.New(cfg)

	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test.toml")

	tomlContent := `
name = "MyApp"
version = "2.0.0"
debug = true
`

	err := os.WriteFile(configFile, []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test conf file: %v", err)
	}

	// 加载配置
	err = manager.LoadFile(configFile, false)
	if err != nil {
		t.Fatalf("Failed to load conf: %v", err)
	}

	// 验证配置
	loadedConfig := manager.GetConfig()
	if loadedConfig.Name != "MyApp" {
		t.Errorf("Expected name 'MyApp', got '%s'", loadedConfig.Name)
	}
	if loadedConfig.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", loadedConfig.Version)
	}
	if !loadedConfig.Debug {
		t.Error("Expected debug to be true")
	}
}

func TestConfigManager_Defaults(t *testing.T) {
	cfg := &TestAppConfig{}
	manager := conf.New(cfg)

	// 创建空配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "empty.toml")

	err := os.WriteFile(configFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty conf file: %v", err)
	}

	// 加载配置
	err = manager.LoadFile(configFile, false)
	if err != nil {
		t.Fatalf("Failed to load conf: %v", err)
	}

	// 验证默认值
	loadedConfig := manager.GetConfig()
	if loadedConfig.Name != "TestApp" {
		t.Errorf("Expected default name 'TestApp', got '%s'", loadedConfig.Name)
	}
	if loadedConfig.Version != "1.0.0" {
		t.Errorf("Expected default version '1.0.0', got '%s'", loadedConfig.Version)
	}
}

func TestConfigManager_UpdateConfig(t *testing.T) {
	cfg := &TestAppConfig{Name: "Original", Version: "1.0.0"}
	manager := conf.New(cfg)

	// 更新配置
	manager.UpdateConfig(func(c *TestAppConfig) {
		c.Name = "Updated"
		c.Debug = true
	})

	// 验证更新
	updatedConfig := manager.GetConfig()
	if updatedConfig.Name != "Updated" {
		t.Errorf("Expected updated name 'Updated', got '%s'", updatedConfig.Name)
	}
	if !updatedConfig.Debug {
		t.Error("Expected debug to be true after update")
	}
	if updatedConfig.Version != "1.0.0" {
		t.Errorf("Expected version to remain '1.0.0', got '%s'", updatedConfig.Version)
	}
}

func TestConfigManager_Watch(t *testing.T) {
	cfg := &TestAppConfig{}
	manager := conf.New(cfg)

	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "watch_test.toml")

	initialContent := `
name = "Initial"
version = "1.0.0"
`

	err := os.WriteFile(configFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test conf file: %v", err)
	}

	// 加载配置并启用监听
	err = manager.LoadFile(configFile, true)
	if err != nil {
		t.Fatalf("Failed to load conf with watch: %v", err)
	}

	// 验证初始配置
	if manager.GetConfig().Name != "Initial" {
		t.Errorf("Expected initial name 'Initial', got '%s'", manager.GetConfig().Name)
	}

	// 修改配置文件
	updatedContent := `
name = "Updated"
version = "2.0.0"
debug = true
`

	err = os.WriteFile(configFile, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to update test conf file: %v", err)
	}

	// 等待文件监听器更新配置
	time.Sleep(100 * time.Millisecond)

	// 验证配置已更新
	updatedConfig := manager.GetConfig()
	if updatedConfig.Name != "Updated" {
		t.Errorf("Expected updated name 'Updated', got '%s'", updatedConfig.Name)
	}
	if updatedConfig.Version != "2.0.0" {
		t.Errorf("Expected updated version '2.0.0', got '%s'", updatedConfig.Version)
	}
}
