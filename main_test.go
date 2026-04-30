package main

import (
	"os"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	content := `
redis:
  addr: "localhost:6379"
  db: 1
channels:
  target: "out"
  sources:
    - "src1"
    - "src2"
`
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()

	cfg, err := loadConfig(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Redis.Addr != "localhost:6379" {
		t.Errorf("Redis.Addr = %q, want %q", cfg.Redis.Addr, "localhost:6379")
	}
	if cfg.Redis.DB != 1 {
		t.Errorf("Redis.DB = %d, want %d", cfg.Redis.DB, 1)
	}
	if cfg.Channels.Target != "out" {
		t.Errorf("Channels.Target = %q, want %q", cfg.Channels.Target, "out")
	}
	if len(cfg.Channels.Sources) != 2 {
		t.Errorf("len(Channels.Sources) = %d, want %d", len(cfg.Channels.Sources), 2)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(": invalid: yaml: {[}"); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()

	_, err = loadConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
