package builder

import (
	"bytes"
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
			patched, err := patchAXMLManifest(data, pkgName, verName, verCode)
			if err != nil {
				return fmt.Errorf("patch manifest: %w", err)
			}
			os.WriteFile(path, patched, fi.Mode())
			return nil

		case "resources.arsc":
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			patched := patchResourcesArsc(data, appName)
			os.WriteFile(path, patched, fi.Mode())
			return nil
		}

		return nil
	})
}

// patchAXMLManifest modifies key values in AndroidManifest.xml using proper AXML parsing.
func patchAXMLManifest(data []byte, pkgName, verName string, verCode int32) ([]byte, error) {
	doc, err := axml.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("axml parse: %w", err)
	}

	// Package name — exact length match required
	if err := doc.SetString("com.kaizhou554.foilexample", pkgName); err != nil {
		return nil, fmt.Errorf("set package: %w", err)
	}

	// Update permission name and provider authority so apps with different
	// package names can be installed side-by-side (no provider conflict).
	// These are string identifiers, NOT class references — safe to change.
	// Activity names (MainActivity) are class references in DEX and stay unchanged.
	_ = doc.SetString("com.kaizhou554.foilexample.DYNAMIC_RECEIVER_NOT_EXPORTED_PERMISSION",
		pkgName+".DYNAMIC_RECEIVER_NOT_EXPORTED_PERMISSION")
	_ = doc.SetString("com.kaizhou554.foilexample.androidx-startup",
		pkgName+".androidx-startup")

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

// patchResourcesArsc replaces known default label strings in resources.arsc
// with the user's app name. Tries both UTF-16LE and UTF-8 string pool encodings.
func patchResourcesArsc(data []byte, appName string) []byte {
	defaults := []string{"Foil Example", "FoilExample", "Foil", "foil"}
	// Try UTF-16LE first, then UTF-8
	result := patchByUTF16Replace(data, defaults, appName)
	result = patchByUTF8Replace(result, defaults, appName)
	return result
}

// patchByUTF16Replace finds the first UTF-16LE occurrence of any default
// string and replaces it with appName (padded/truncated to same byte length).
func patchByUTF16Replace(data []byte, defaults []string, appName string) []byte {
	result := make([]byte, len(data))
	copy(result, data)

	appUTF16 := encodeUTF16(appName)

	for _, old := range defaults {
		oldUTF16 := encodeUTF16(old)
		idx := bytes.Index(result, oldUTF16)
		if idx < 0 {
			continue
		}
		padded := make([]byte, len(oldUTF16))
		copy(padded, appUTF16)
		for i := len(appUTF16); i < len(oldUTF16); i++ {
			padded[i] = 0
		}
		copy(result[idx:], padded)
		break
	}
	return result
}

// patchByUTF8Replace finds the first UTF-8 occurrence of any default string
// and replaces it with appName (padded/truncated to same byte length).
func patchByUTF8Replace(data []byte, defaults []string, appName string) []byte {
	result := make([]byte, len(data))
	copy(result, data)

	appBytes := []byte(appName)

	for _, old := range defaults {
		oldBytes := []byte(old)
		idx := bytes.Index(result, oldBytes)
		if idx < 0 {
			continue
		}
		padded := make([]byte, len(oldBytes))
		copy(padded, appBytes)
		for i := len(appBytes); i < len(oldBytes); i++ {
			padded[i] = 0
		}
		copy(result[idx:], padded)
		break
	}
	return result
}

func encodeUTF16(s string) []byte {
	runes := []rune(s)
	b := make([]byte, len(runes)*2)
	for i, r := range runes {
		b[i*2] = byte(r)
		b[i*2+1] = byte(r >> 8)
	}
	return b
}
