package builder

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


// resolveJavaPath finds the java executable, trying both platform variants.
func resolveJavaPath(assetDir string) string {
	base := filepath.Join(assetDir, "jre-minimal", "bin", "java")
	if _, err := os.Stat(base); err == nil {
		return base
	}
	// Windows: try with .exe suffix
	exe := base + ".exe"
	if _, err := os.Stat(exe); err == nil {
		return exe
	}
	return base // will fail later with a clear error
}

// apktoolPatchManifest uses apktool to decode the template APK, applies
// manifest and resource-string modifications in plain text, rebuilds the APK,
// and extracts the modified binary files (AndroidManifest.xml + resources.arsc)
// into unpackDir, overwriting the originals.
//
// This replaces the pure-Go AXML/ARSC patching in patch.go.
func apktoolPatchManifest(unpackDir, templatePath, pkgName, appName, verName string, verCode int32) error {
	// Resolve java + jar paths relative to the template APK's directory.
	// templatePath is typically "assets/foil-example.apk", so the asset
	// directory is the parent of the template.
	assetDir := filepath.Dir(templatePath)
	javaPath := resolveJavaPath(assetDir)
	jarPath := filepath.Join(assetDir, "apktool.jar")

	// Verify assets exist
	if _, err := os.Stat(javaPath); err != nil {
		return fmt.Errorf("bundled java not found at %s: %w", javaPath, err)
	}
	if _, err := os.Stat(jarPath); err != nil {
		return fmt.Errorf("apktool jar not found at %s: %w", jarPath, err)
	}

	// Create temp workspace
	tmpDir, err := os.MkdirTemp("", "apktool-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	decodeDir := filepath.Join(tmpDir, "decoded")
	rebuiltPath := filepath.Join(tmpDir, "rebuilt.apk")

	// ── Step 1: Decode ──────────────────────────────────────────────────
	// --no-src skips DEX decompilation (not needed, much faster)
	out, err := exec.Command(javaPath, "-jar", jarPath,
		"d", templatePath, "-o", decodeDir, "--no-src").CombinedOutput()
	if err != nil {
		return fmt.Errorf("apktool decode: %w\n%s", err, string(out))
	}

	// ── Step 2: Modify text AndroidManifest.xml ─────────────────────────
	if err := modifyTextManifest(
		filepath.Join(decodeDir, "AndroidManifest.xml"),
		pkgName, verName, verCode); err != nil {
		return fmt.Errorf("modify manifest: %w", err)
	}

	// ── Step 3: Update app-name string in decoded resources ────────────
	if appName != "" {
		if err := updateAppNameString(decodeDir, appName); err != nil {
			// Non-fatal: the label might be set differently in some APKs
			fmt.Printf("[apktool] warning: update app name: %v\n", err)
		}
	}

	// ── Step 4: Rebuild ────────────────────────────────────────────────
	out, err = exec.Command(javaPath, "-jar", jarPath,
		"b", decodeDir, "-o", rebuiltPath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("apktool build: %w\n%s", err, string(out))
	}

	// ── Step 5: Extract modified binary files → overwrite in unpackDir ─
	for _, name := range []string{"AndroidManifest.xml", "resources.arsc"} {
		dest := filepath.Join(unpackDir, name)
		if err := extractFromZip(rebuiltPath, name, dest); err != nil {
			return fmt.Errorf("extract %s: %w", name, err)
		}
	}

	return nil
}

// ── Text manifest editing ──────────────────────────────────────────────────

const (
	oldPkg = "com.kaizhou554.foilexample"
)

// modifyTextManifest applies all manifest changes to the decoded text XML.
func modifyTextManifest(manifestPath, pkgName, verName string, verCode int32) error {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	content := string(data)

	// 1. Package name
	content = strings.ReplaceAll(content,
		`package="`+oldPkg+`"`,
		`package="`+pkgName+`"`)

	// 2. Permission name
	oldPerm := oldPkg + ".DYNAMIC_RECEIVER_NOT_EXPORTED_PERMISSION"
	newPerm := pkgName + ".DYNAMIC_RECEIVER_NOT_EXPORTED_PERMISSION"
	content = strings.ReplaceAll(content, oldPerm, newPerm)

	// 3. Provider authority
	oldAuth := oldPkg + ".androidx-startup"
	newAuth := pkgName + ".androidx-startup"
	content = strings.ReplaceAll(content, oldAuth, newAuth)

	// 4. Version attributes on the <manifest> root element
	//    Add them if missing; update if present.
	content = upsertManifestAttr(content, "android:versionName", verName)
	content = upsertManifestAttr(content, "android:versionCode", fmt.Sprintf("%d", verCode))

	return os.WriteFile(manifestPath, []byte(content), 0644)
}

// upsertManifestAttr looks for attr="value" inside the <manifest ...> opening
// tag.  If found it replaces the value; if not it inserts the attribute before
// the closing '>' of the <manifest> tag.
func upsertManifestAttr(xml, attr, value string) string {
	// Find the <manifest ...> tag first (skip <?xml ... ?>)
	manifestStart := strings.Index(xml, "<manifest")
	if manifestStart < 0 {
		return xml
	}

	// Search inside the <manifest ...> region (until the closing >)
	tagSection := xml[manifestStart:]

	marker := attr + `="`
	beforeRel := strings.Index(tagSection, marker)
	if beforeRel >= 0 {
		// Found — replace value between the quotes
		before := manifestStart + beforeRel
		quoteStart := before + len(marker)
		quoteEnd := strings.IndexByte(xml[quoteStart:], '"')
		if quoteEnd >= 0 {
			return xml[:quoteStart] + value + xml[quoteStart+quoteEnd:]
		}
	}

	// Not found — find the closing > of the <manifest> opening tag
	tagEndRel := strings.IndexByte(tagSection, '>')
	if tagEndRel < 0 {
		return xml
	}
	insert := fmt.Sprintf(` %s="%s"`, attr, value)
	tagEnd := manifestStart + tagEndRel
	return xml[:tagEnd] + insert + xml[tagEnd:]
}

// ── Resource-string editing ────────────────────────────────────────────────

// updateAppNameString finds the default-locale strings.xml in the decoded
// resources and updates the <string name="app_name">…</string> value.
func updateAppNameString(decodeDir, appName string) error {
	// The default locale strings live in res/values/strings.xml (no -xx suffix)
	stringsPath := filepath.Join(decodeDir, "res", "values", "strings.xml")
	data, err := os.ReadFile(stringsPath)
	if err != nil {
		return fmt.Errorf("read strings.xml: %w", err)
	}

	content := string(data)

	// Look for <string name="app_name">...</string>
	marker := `<string name="app_name">`
	start := strings.Index(content, marker)
	if start < 0 {
		return fmt.Errorf("app_name string not found in strings.xml")
	}
	start += len(marker)
	end := strings.Index(content[start:], "</string>")
	if end < 0 {
		return fmt.Errorf("app_name closing tag not found")
	}

	// Escape XML special characters in the new name
	escaped := escapeXML(appName)

	content = content[:start] + escaped + content[start+end:]
	return os.WriteFile(stringsPath, []byte(content), 0644)
}

// escapeXML escapes the five characters that have special meaning in XML text.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

// ── ZIP extraction ─────────────────────────────────────────────────────────

// extractFromZip extracts a single file from a ZIP archive to destPath.
func extractFromZip(zipPath, name, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		os.MkdirAll(filepath.Dir(destPath), 0755)
		out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = copyFile(out, rc)
		if err != nil {
			return fmt.Errorf("write %s: %w", destPath, err)
		}
		return nil
	}

	return fmt.Errorf("entry %q not found in %s", name, zipPath)
}

// copyFile copies from reader to writer (standard io copy).
func copyFile(dst *os.File, src io.ReadCloser) (int64, error) {
	return io.Copy(dst, src)
}
