package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

// Manager handles loading, saving, and thread-safe access to configuration.
type Manager struct {
	path string
	data *Config
	mu   sync.RWMutex
}

// NewManager creates a new Manager for the config file at the given path.
// It automatically creates parent directories and a default config file if they don't exist.
// Errors are returned if directory creation or initial config write fails.
func NewManager(configPath string) (*Manager, error) {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create config directory %s: %w", dir, err)
	}

	m := &Manager{
		path: configPath,
	}

	// If the config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		m.data = DefaultConfig()
		if err := m.saveUnsafe(); err != nil {
			return nil, fmt.Errorf("write default config: %w", err)
		}
		return m, nil
	}

	// Load existing config; on any failure, back up the corrupt file
	// and regenerate with defaults so the app can still start.
	if err := m.loadUnsafe(); err != nil {
		return m.recoverCorrupted(err)
	}
	if err := m.data.Validate(); err != nil {
		return m.recoverCorrupted(err)
	}

	return m, nil
}

// Load reads the TOML config file from disk.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.loadUnsafe()
}

// Save writes the current configuration to disk, overwriting the old file.
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.saveUnsafe()
}

// Get returns a copy of the current configuration.
// Safe for concurrent use.
func (m *Manager) Get() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Return a defensive copy so callers can't mutate the internal state
	cpy := *m.data
	return &cpy
}

// Update applies fn to the current configuration and persists the result.
// fn receives a mutable pointer; any changes are saved to disk after fn returns.
// If fn returns an error, the config is NOT saved and the error is propagated.
func (m *Manager) Update(fn func(*Config) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := fn(m.data); err != nil {
		return err
	}

	return m.saveUnsafe()
}

// ── Recovery ──

// recoverCorrupted renames the broken config file to a timestamped backup
// and regenerates it with defaults so the application can still start.
func (m *Manager) recoverCorrupted(cause error) (*Manager, error) {
	backup := m.path + ".corrupted." + time.Now().Format("20060102-150405")
	if err := os.Rename(m.path, backup); err != nil {
		return nil, fmt.Errorf("config corrupted (%w) and failed to backup bad file: %w", cause, err)
	}
	println("Warning: config file was corrupted (" + cause.Error() + "), backed up to " + backup + " and reset to defaults")
	m.data = DefaultConfig()
	if err := m.saveUnsafe(); err != nil {
		return nil, fmt.Errorf("config corrupted (%w), backed up to %s, but failed to write fresh defaults: %w", cause, backup, err)
	}
	return m, nil
}

// ── Internal (caller must hold the appropriate lock) ──

func (m *Manager) loadUnsafe() error {
	m.data = DefaultConfig()
	if _, err := toml.DecodeFile(m.path, m.data); err != nil {
		return fmt.Errorf("decode config file %s: %w", m.path, err)
	}
	return nil
}

func (m *Manager) saveUnsafe() error {
	f, err := os.Create(m.path)
	if err != nil {
		return fmt.Errorf("open config file %s for writing: %w", m.path, err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(m.data); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}
