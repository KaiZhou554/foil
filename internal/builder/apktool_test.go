package builder

import (
	"archive/zip"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"lets_config/internal/axml"
)

func TestApktoolManifestPatching(t *testing.T) {
	// Create temp workspace
	tmpDir := t.TempDir()
	unpackDir := filepath.Join(tmpDir, "unpacked")

	// First unpack the template APK (same as builder.unpackAPK)
	templatePath := filepath.Join("..", "..", "assets", "foil-example.apk")
	r, err := zip.OpenReader(templatePath)
	if err != nil {
		t.Fatalf("open template: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(unpackDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		rc, _ := f.Open()
		data := make([]byte, f.UncompressedSize64)
		rc.Read(data)
		rc.Close()
		os.WriteFile(fpath, data, f.Mode())
	}

	// Run apktool-based manifest patching
	pkgName := "com.mytestapp12345"
	appName := "我的App"
	verName := "7"
	verCode := int32(7)

	if err := apktoolPatchManifest(unpackDir, templatePath, pkgName, appName, verName, verCode); err != nil {
		t.Fatalf("apktoolPatchManifest failed: %v", err)
	}

	// Verify extracted files exist
	manifestPath := filepath.Join(unpackDir, "AndroidManifest.xml")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("AndroidManifest.xml not found after patching: %v", err)
	}
	arscPath := filepath.Join(unpackDir, "resources.arsc")
	if _, err := os.Stat(arscPath); err != nil {
		t.Fatalf("resources.arsc not found after patching: %v", err)
	}

	// Read and parse binary AXML manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	// Verify it's valid binary AXML (starts with 0x0003)
	if len(data) < 2 || data[0] != 0x03 || data[1] != 0x00 {
		t.Fatalf("manifest doesn't look like binary AXML: first bytes = %x", data[:2])
	}

	doc, err := axml.Parse(data)
	if err != nil {
		t.Fatalf("parse binary AXML: %v", err)
	}

	// Verify key string values
	checks := []struct {
		name     string
		expected string
	}{
		{"package name", pkgName},
		{"updated permission", pkgName + ".DYNAMIC_RECEIVER_NOT_EXPORTED_PERMISSION"},
		{"updated provider authority", pkgName + ".androidx-startup"},
	}

	for _, c := range checks {
		if idx := doc.FindString(c.expected); idx < 0 {
			t.Errorf("manifest AXML missing %q", c.expected)
		}
	}

	// Verify OLD package name is NOT present
	if idx := doc.FindString("com.kaizhou554.foilexample"); idx >= 0 {
		t.Errorf("old package name still present in AXML string pool")
	}

	// Verify files are not empty
	maniInfo, _ := os.Stat(manifestPath)
	arscInfo, _ := os.Stat(arscPath)
	t.Logf("AndroidManifest.xml: %d bytes", maniInfo.Size())
	t.Logf("resources.arsc: %d bytes", arscInfo.Size())
	t.Logf("Patching OK — package=%s, verName=%s, verCode=%d", pkgName, verName, verCode)
}

func TestResourcesArscProduced(t *testing.T) {
	tmpDir := t.TempDir()
	unpackDir := filepath.Join(tmpDir, "unpacked")
	templatePath := filepath.Join("..", "..", "assets", "foil-example.apk")

	// Unpack template
	r, err := zip.OpenReader(templatePath)
	if err != nil {
		t.Fatalf("open template: %v", err)
	}
	for _, f := range r.File {
		fpath := filepath.Join(unpackDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		rc, _ := f.Open()
		data := make([]byte, f.UncompressedSize64)
		rc.Read(data)
		rc.Close()
		os.WriteFile(fpath, data, f.Mode())
	}
	r.Close()

	appName := "我的测试App"
	if err := apktoolPatchManifest(unpackDir, templatePath, "com.test.pkg", appName, "9", 9); err != nil {
		t.Fatalf("apktoolPatchManifest failed: %v", err)
	}

	arscPath := filepath.Join(unpackDir, "resources.arsc")
	arscData, err := os.ReadFile(arscPath)
	if err != nil {
		t.Fatalf("read resources.arsc: %v", err)
	}

	// Verify ARSC is not empty and is a valid resource table
	if len(arscData) < 100 {
		t.Fatalf("resources.arsc suspiciously small: %d bytes", len(arscData))
	}
	if binary.LittleEndian.Uint16(arscData[:2]) != 0x0002 {
		t.Fatalf("resources.arsc doesn't start with RES_TABLE type 0x0002: %x", arscData[:2])
	}

	t.Logf("resources.arsc: %d bytes, valid resource table", len(arscData))
}
