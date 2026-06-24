package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lets_config/internal/axml"
)

// patchAllStrings walks unpackDir and patches AndroidManifest.xml + resources.arsc.
// Uses proper AXML parsing — no raw byte replacement that could cause corruption.
func patchAllStrings(unpackDir, pkgName, appName, verName string, verCode int32) error {
	return filepath.Walk(unpackDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}

		base := strings.ToLower(filepath.Base(path))

		switch base {
		case "androidmanifest.xml":
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			patched, err := patchAXMLManifest(data, pkgName, appName, verName, verCode)
			if err != nil {
				return fmt.Errorf("patch manifest: %w", err)
			}
			os.WriteFile(path, patched, fi.Mode())
			return nil

		case "resources.arsc":
			// TODO: patch resources.arsc for app label
			return nil
		}

		return nil
	})
}

// patchAXMLManifest modifies key values in AndroidManifest.xml using proper AXML parsing.
func patchAXMLManifest(data []byte, pkgName, appName, verName string, verCode int32) ([]byte, error) {
	doc, err := axml.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("axml parse: %w", err)
	}

	// Package name — exact length match required
	if err := doc.SetString("com.kaizhou554.foilexample", pkgName); err != nil {
		return nil, fmt.Errorf("set package: %w", err)
	}

	// NOTE: We ONLY change the manifest's `package` attribute.
	// Activity/permission class names like .MainActivity and
	// .DYNAMIC_RECEIVER... are compiled into classes.dex and CANNOT
	// be changed without recompiling the Java source. They reference
	// the original package and Android resolves them correctly as
	// long as they're fully qualified in the manifest.

	// Version name — exact length match
	if err := doc.SetString("1.0", verName); err != nil {
		// Try alternative version formats that might be in the pool
		_ = err
	}

	// Update versionCode (integer attribute, not string)
	if verCode > 0 {
		if err := doc.SetIntAttribute("versionCode", verCode); err != nil {
			_ = err // non-fatal
		}
	}

	return doc.Raw(), nil
}
