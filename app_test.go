package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareFileInputZipSlip(t *testing.T) {
	base := TempDirBase()
	os.MkdirAll(base, 0755)
	defer os.RemoveAll(base)

	// Create a malicious ZIP with path traversal
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "evil.zip")

	// Build a ZIP containing a path traversal entry
	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(zf)

	// Add a legitimate index.html
	w, _ := zw.Create("index.html")
	w.Write([]byte("<html></html>"))

	// Add a path traversal entry
	w2, _ := zw.Create("../../../../Windows/System32/evil.dll")
	w2.Write([]byte("malicious"))

	zw.Close()
	zf.Close()

	// Try to extract this ZIP — should fail with a path traversal error
	app := &App{}
	_, err = app.PrepareFileInput(zipPath)
	if err == nil {
		t.Error("expected error for ZIP with path traversal, got nil")
	}
	if !strings.Contains(err.Error(), "tries to escape") {
		t.Errorf("expected path traversal error, got: %v", err)
	}
}

func TestPrepareFileInputZipSlipAbsolute(t *testing.T) {
	base := TempDirBase()
	os.MkdirAll(base, 0755)
	defer os.RemoveAll(base)

	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "evil2.zip")

	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(zf)

	// Add index.html
	w, _ := zw.Create("index.html")
	w.Write([]byte("<html></html>"))

	// Add absolute path entry (on Windows: C:\Windows\...)
	w2, _ := zw.Create(`C:\Windows\System32\evil.dll`)
	w2.Write([]byte("malicious"))

	zw.Close()
	zf.Close()

	app := &App{}
	_, err = app.PrepareFileInput(zipPath)
	if err == nil {
		t.Error("expected error for ZIP with absolute path, got nil")
	}
}

func TestPrepareFileInputEmptyPath(t *testing.T) {
	app := &App{}
	_, err := app.PrepareFileInput("")
	if err == nil {
		t.Error("expected error for empty file path")
	}
}

func TestPrepareFileInputNonexistentFile(t *testing.T) {
	app := &App{}
	_, err := app.PrepareFileInput(`C:\nonexistent\file.html`)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestPrepareFileInputUnsupportedType(t *testing.T) {
	app := &App{}
	_, err := app.PrepareFileInput(`C:\test.txt`)
	if err == nil {
		t.Error("expected error for unsupported file type")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected 'unsupported' error, got: %v", err)
	}
}

func TestPrepareFileInputHtmlFile(t *testing.T) {
	base := TempDirBase()
	os.MkdirAll(base, 0755)
	defer os.RemoveAll(base)

	// Create a temporary HTML file
	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(htmlPath, []byte("<html><body>Test</body></html>"), 0644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	app := &App{}
	projectDir, err := app.PrepareFileInput(htmlPath)
	if err != nil {
		t.Fatalf("PrepareFileInput failed: %v", err)
	}

	// Verify index.html was created in the project dir
	indexPath := filepath.Join(projectDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("index.html not found in project dir: %v", err)
	}

	// Cleanup
	os.RemoveAll(projectDir)
}

func TestPrepareFileInputValidZip(t *testing.T) {
	base := TempDirBase()
	os.MkdirAll(base, 0755)
	defer os.RemoveAll(base)

	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "project.zip")

	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(zf)

	w, _ := zw.Create("index.html")
	w.Write([]byte("<html><body>Test</body></html>"))

	// Also add an assets subdirectory with JS
	w2, _ := zw.Create("assets/app.js")
	w2.Write([]byte("console.log('hello')"))

	zw.Close()
	zf.Close()

	app := &App{}
	projectDir, err := app.PrepareFileInput(zipPath)
	if err != nil {
		t.Fatalf("PrepareFileInput failed: %v", err)
	}

	// Verify index.html is in the resolved project dir
	indexPath := filepath.Join(projectDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("index.html not found: %v", err)
	}

	// Cleanup
	os.RemoveAll(projectDir)
}

func TestAssetsDirDevMode(t *testing.T) {
	// In dev mode (assets/ exists), should return "assets"
	dir := AssetsDir()
	if dir != "assets" {
		t.Logf("AssetsDir in test: %s", dir)
	}
}

func TestTempDirBase(t *testing.T) {
	base := TempDirBase()
	if base == "" {
		t.Error("TempDirBase returned empty string")
	}
	if !strings.Contains(base, "unieditdept") {
		t.Errorf("TempDirBase should contain 'unieditdept', got: %s", base)
	}
	if !strings.Contains(base, "foil") {
		t.Errorf("TempDirBase should contain 'foil', got: %s", base)
	}
}

func TestGetAppVersion(t *testing.T) {
	app := &App{}
	v := app.GetAppVersion()
	if v == "" {
		t.Error("AppVersion should not be empty")
	}
	t.Logf("AppVersion: %s", v)
}

// Fuzz tests for app.go functions

func FuzzPrepareFileInput(f *testing.F) {
	f.Add("test.html")
	f.Add("test.zip")
	f.Add("")
	f.Add("../etc/passwd")
	f.Fuzz(func(t *testing.T, filePath string) {
		app := &App{}
		// Must never panic
		result, _ := app.PrepareFileInput(filePath)
		// Cleanup if result is from temp dir
		if result != "" && strings.HasPrefix(result, TempDirBase()) {
			os.RemoveAll(result)
		}
	})
}
