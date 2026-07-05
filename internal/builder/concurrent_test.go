package builder

import (
	"sync"
	"testing"
)

// TestGeneratePackageNameConcurrent verifies concurrent name generation produces unique names.
func TestGeneratePackageNameConcurrent(t *testing.T) {
	const (
		numGoroutines = 20
		namesPerRoutine = 100
	)

	var (
		mu   sync.Mutex
		seen = make(map[string]bool, numGoroutines*namesPerRoutine)
		wg   sync.WaitGroup
	)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < namesPerRoutine; j++ {
				pkg := GeneratePackageName("TestApp", nil, 26)
				mu.Lock()
				if seen[pkg] {
					t.Errorf("duplicate package name from concurrent generation: %q", pkg)
				}
				seen[pkg] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(seen) != numGoroutines*namesPerRoutine {
		t.Errorf("expected %d unique names, got %d", numGoroutines*namesPerRoutine, len(seen))
	}
}

// TestGeneratePackageNameConcurrentMultiApp tests concurrent generation with different app names.
func TestGeneratePackageNameConcurrentMultiApp(t *testing.T) {
	apps := []string{"App1", "App2", "App3", "我的应用", ""}
	const namesPerApp = 50

	var (
		mu   sync.Mutex
		seen = make(map[string]bool)
		wg   sync.WaitGroup
	)

	for _, app := range apps {
		app := app
		for g := 0; g < 4; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < namesPerApp; i++ {
					pkg := GeneratePackageName(app, nil, 26)
					mu.Lock()
					if seen[pkg] {
						t.Errorf("duplicate: app=%q pkg=%q", app, pkg)
					}
					seen[pkg] = true
					mu.Unlock()
				}
			}()
		}
	}

	wg.Wait()
	t.Logf("Generated %d unique package names across %d goroutines", len(seen), len(apps)*4)
}

// Benchmark for concurrent package name generation
func BenchmarkGeneratePackageNameParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GeneratePackageName("BenchApp", nil, 26)
		}
	})
}
