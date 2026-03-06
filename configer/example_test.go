package configer_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bizvip/go-utils/configer"
)

type exampleConfig struct {
	Name  string `json:"name"`
	Debug bool   `json:"debug"`
}

func (c *exampleConfig) SetDefaults() {
	if c.Name == "" {
		c.Name = "default"
	}
}

func (c *exampleConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func jsonDecoder(data []byte, out any) error {
	return json.Unmarshal(data, out)
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(path, []byte(`{"name":"demo","debug":true}`), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	cfg, err := configer.Load[exampleConfig](path,
		configer.WithDecoder[exampleConfig](configer.DecoderFunc(jsonDecoder)),
	)
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}
	if cfg.Name != "demo" || !cfg.Debug {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestManagerWatch(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(path, []byte(`{"name":"initial"}`), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	manager := configer.NewManager[exampleConfig](path,
		configer.WithDecoder[exampleConfig](configer.DecoderFunc(jsonDecoder)),
	)
	if _, err := manager.Load(); err != nil {
		t.Fatalf("initial load failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	changed := make(chan string, 1)
	if err := manager.Watch(ctx, func(cfg *exampleConfig) {
		changed <- cfg.Name
	}, func(err error) {
		t.Errorf("watch error: %v", err)
	}); err != nil {
		t.Fatalf("watch failed: %v", err)
	}

	if err := os.WriteFile(path, []byte(`{"name":"updated"}`), 0o644); err != nil {
		t.Fatalf("rewrite config failed: %v", err)
	}

	select {
	case value := <-changed:
		if value != "updated" {
			t.Fatalf("unexpected watched value: %s", value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watch timeout")
	}
}
