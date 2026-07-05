package config

import (
	"path/filepath"
	"sync"
	"testing"
)

// TestManagerConcurrentAccess verifies the Manager handles concurrent reads and writes.
func TestManagerConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	m, err := NewManager(filepath.Join(dir, "config.toml"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	const (
		numReaders = 50
		numWriters = 10
		iterations = 100
	)

	var wg sync.WaitGroup

	// Start concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				cfg := m.Get()
				if cfg == nil {
					t.Error("Get() returned nil")
					return
				}
				if cfg.Version == "" {
					t.Error("Version should not be empty")
					return
				}
			}
		}()
	}

	// Start concurrent writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				err := m.Update(func(c *Config) error {
					c.Version = "1.0"
					c.Language = "en"
					return nil
				})
				if err != nil {
					t.Errorf("Update failed: %v", err)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify final state is consistent
	cfg := m.Get()
	if cfg.Version != "1.0" {
		t.Errorf("final Version = %q, want %q", cfg.Version, "1.0")
	}
}

// TestManagerConcurrentReadWrite verifies concurrent reads and writes don't cause data races.
func TestManagerConcurrentReadWrite(t *testing.T) {
	dir := t.TempDir()
	m, err := NewManager(filepath.Join(dir, "config.toml"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	var wg sync.WaitGroup

	// Rapid concurrent read+write
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				_ = m.Get()
			} else {
				_ = m.Update(func(c *Config) error {
					c.Language = "en"
					return nil
				})
			}
		}(i)
	}

	wg.Wait()
}

// TestManagerUpdateRollbackConcurrent verifies failed updates don't corrupt state under concurrency.
func TestManagerUpdateRollbackConcurrent(t *testing.T) {
	dir := t.TempDir()
	m, err := NewManager(filepath.Join(dir, "config.toml"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	// Set initial state
	if err := m.Update(func(c *Config) error {
		c.Language = "zh-CN"
		return nil
	}); err != nil {
		t.Fatalf("initial Update: %v", err)
	}

	var wg sync.WaitGroup
	errCount := 0
	var errMu sync.Mutex

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := m.Update(func(c *Config) error {
				c.Language = "en"
				return testSentinelError{}
			})
			if err != nil {
				errMu.Lock()
				errCount++
				errMu.Unlock()
			}
		}()
	}

	wg.Wait()

	// All updates should have failed
	if errCount != 50 {
		t.Errorf("expected 50 failed updates, got %d", errCount)
	}

	// Original state must be preserved
	cfg := m.Get()
	if cfg.Language != "zh-CN" {
		t.Errorf("Language should be 'zh-CN', got %q", cfg.Language)
	}
}

// testSentinelError is a simple error type for testing rollback.
type testSentinelError struct{}

func (e testSentinelError) Error() string { return "test sentinel error" }

// TestManagerHighConcurrencyGet stress-tests Get() under high concurrency.
func TestManagerHighConcurrencyGet(t *testing.T) {
	dir := t.TempDir()
	m, err := NewManager(filepath.Join(dir, "config.toml"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = m.Get()
			}
		}()
	}
	wg.Wait()
}

// Benchmarks for concurrent access

func BenchmarkManagerGetParallel(b *testing.B) {
	dir := b.TempDir()
	m, _ := NewManager(filepath.Join(dir, "config.toml"))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = m.Get()
		}
	})
}

func BenchmarkManagerUpdateParallel(b *testing.B) {
	dir := b.TempDir()
	m, _ := NewManager(filepath.Join(dir, "config.toml"))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = m.Update(func(c *Config) error {
				c.Version = "1.0"
				return nil
			})
		}
	})
}
