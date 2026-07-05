package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Version == "" {
		t.Error("Version must not be empty")
	}
	if cfg.Language == "" {
		t.Error("Language must not be empty")
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"valid", Config{Version: "1.0", Language: "en"}, false},
		{"empty version", Config{Version: "", Language: "en"}, true},
		{"empty language", Config{Version: "1.0", Language: ""}, true},
		{"both empty", Config{Version: "", Language: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewManagerCreatesDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	cfg := m.Get()
	if cfg == nil {
		t.Fatal("Get() returned nil")
	}
	if cfg.Version == "" {
		t.Error("Version should not be empty")
	}
	if !cfg.FirstLaunch {
		t.Error("FirstLaunch should be true by default")
	}

	// Verify file was created
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

func TestManagerLoadsExistingConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	// Create first
	m1, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager (create): %v", err)
	}
	if err := m1.Update(func(c *Config) error {
		c.Language = "en"
		c.FirstLaunch = false
		return nil
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	// Load again
	m2, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager (load): %v", err)
	}

	cfg := m2.Get()
	if cfg.Language != "en" {
		t.Errorf("Language = %q, want %q", cfg.Language, "en")
	}
	if cfg.FirstLaunch {
		t.Error("FirstLaunch should be false")
	}
}

func TestManagerUpdateErrorRollback(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	// Update that succeeds
	if err := m.Update(func(c *Config) error {
		c.Language = "en"
		return nil
	}); err != nil {
		t.Fatalf("first update: %v", err)
	}

	// Update that fails should NOT change config
	err = m.Update(func(c *Config) error {
		c.Language = "fr"
		return os.ErrPermission // simulate error
	})
	if err == nil {
		t.Fatal("expected error from Update")
	}

	cfg := m.Get()
	if cfg.Language != "en" {
		t.Errorf("Language should remain 'en', got %q", cfg.Language)
	}
}

func TestManagerGetReturnsCopy(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	cfg := m.Get()
	cfg.Language = "modified"

	// Original should be unchanged
	cfg2 := m.Get()
	if cfg2.Language == "modified" {
		t.Error("Get() should return a copy, not a reference")
	}
}

func TestManagerRecoversCorruptedConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	// Create a valid config first
	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if err := m.Update(func(c *Config) error {
		c.Language = "en"
		return nil
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	// Corrupt the file manually
	if err := os.WriteFile(configPath, []byte("this is not valid toml {{{"), 0644); err != nil {
		t.Fatalf("write corrupt config: %v", err)
	}

	// Re-open: should recover and create backup
	m2, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager should recover: %v", err)
	}

	cfg := m2.Get()
	if cfg == nil {
		t.Fatal("config should not be nil after recovery")
	}

	// Verify backup was created
	backups, _ := filepath.Glob(configPath + ".corrupted.*")
	if len(backups) == 0 {
		t.Error("no backup file created for corrupted config")
	}
	t.Logf("Backup created: %v", backups)
}

func TestConfigDirectoryCreation(t *testing.T) {
	dir := t.TempDir()
	nestedPath := filepath.Join(dir, "deeply", "nested", "dir", "config.toml")

	m, err := NewManager(nestedPath)
	if err != nil {
		t.Fatalf("NewManager with nested path: %v", err)
	}

	cfg := m.Get()
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
}

func TestConfigWithSpecialChars(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	// Test with Unicode app name and other special fields
	if err := m.Update(func(c *Config) error {
		c.CompanyName = "测试公司™"
		c.OutputDir = filepath.Join(dir, "输出目录", "sub dir")
		return nil
	}); err != nil {
		t.Fatalf("Update with Unicode: %v", err)
	}

	cfg := m.Get()
	if cfg.CompanyName != "测试公司™" {
		t.Errorf("CompanyName = %q, want %q", cfg.CompanyName, "测试公司™")
	}
}

func TestConfigVeryLongValues(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	m, err := NewManager(configPath)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	longStr := strings.Repeat("x", 10000)
	if err := m.Update(func(c *Config) error {
		c.CompanyName = longStr
		return nil
	}); err != nil {
		t.Fatalf("Update with long value: %v", err)
	}

	cfg := m.Get()
	if cfg.CompanyName != longStr {
		t.Errorf("long string corrupted: expected %d chars, got %d", len(longStr), len(cfg.CompanyName))
	}
}

func TestConfigEmptyPath(t *testing.T) {
	_, err := NewManager("")
	if err == nil {
		t.Error("expected error for empty config path")
	}
}

func BenchmarkNewManager(b *testing.B) {
	dir := b.TempDir()
	for i := 0; i < b.N; i++ {
		path := filepath.Join(dir, "bench", "config.toml")
		os.RemoveAll(filepath.Dir(path))
		_, _ = NewManager(path)
	}
}

func BenchmarkManagerGet(b *testing.B) {
	dir := b.TempDir()
	m, _ := NewManager(filepath.Join(dir, "config.toml"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Get()
	}
}

func BenchmarkManagerUpdate(b *testing.B) {
	dir := b.TempDir()
	m, _ := NewManager(filepath.Join(dir, "config.toml"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Update(func(c *Config) error {
			c.Version = "1.0"
			return nil
		})
	}
}
