package builder

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"lets_config/internal/apksigner/android"
	"lets_config/internal/apksigner/apksign"
)

// ── Public types ──────────────────────────────────────────────────────────

// BuildInput describes what the user wants to build.
type BuildInput struct {
	ProjectDir   string            // path to HTML project (must contain index.html)
	AppName      string            // user-visible app name
	PackageName  string            // optional; empty = auto-generate
	IconWebP     map[string][]byte // path -> webp bytes, e.g. "mipmap-hdpi-v4/ic_launcher.webp" -> data
}

// BuildResult describes what was produced.
type BuildResult struct {
	APKPath     string
	PackageName string
	VersionName string
	VersionCode int32
	Log         string
}

// Builder orchestrates the APK build pipeline.
type Builder struct {
	TemplatePath string // path to foil-example.apk
	OutputDir    string // where final APK goes
	KeysDir      string // where signing keys live
	WorkDir      string // temp workspace

	KeepWorkDir bool // set true to keep temp files for debugging

	logBuf bytes.Buffer
	logger *log.Logger
}

// New creates a Builder with default paths.
func New(templatePath, outputDir, keysDir, workDir string) *Builder {
	b := &Builder{
		TemplatePath: templatePath,
		OutputDir:    outputDir,
		KeysDir:      keysDir,
		WorkDir:      workDir,
	}
	b.logger = log.New(&b.logBuf, "", log.Ltime|log.Lshortfile)
	return b
}

func (b *Builder) logf(format string, args ...interface{}) {
	b.logger.Printf(format, args...)
	fmt.Printf("[builder] "+format+"\n", args...)
}

// ── Pipeline steps ────────────────────────────────────────────────────────

// Build runs the full pipeline end-to-end.
func (b *Builder) Build(in BuildInput) (result *BuildResult, err error) {
	// Catch panics and convert to errors for better debugging
	defer func() {
		if r := recover(); r != nil {
			// Get stack trace
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			err = fmt.Errorf("PANIC: %v\nStack:\n%s", r, buf[:n])
		}
	}()

	b.logBuf.Reset()
	b.logf("=== Build started ===")
	start := time.Now()

	// 1. Validate input
	if err := b.validateInput(in); err != nil {
		return nil, b.fail("validation", err)
	}

	// 2. Create build workspace
	taskDir, err := b.createWorkspace()
	if err != nil {
		return nil, b.fail("workspace", err)
	}

	// 3. Resolve or auto-generate metadata
	const oldPkgLen = 26 // "com.kaizhou554.foilexample" is exactly 26 chars
	pkgName := in.PackageName
	if pkgName == "" {
		pkgName = GeneratePackageName(in.AppName, nil, oldPkgLen)
	}
	verName := GenerateVersionName()
	verCode := GenerateVersionCode()
	b.logf("Package: %s (%d chars) | Version: %s (%d)", pkgName, len(pkgName), verName, verCode)

	// 4. Unpack template APK
	unpackDir := filepath.Join(taskDir, "unpacked")
	if err := b.unpackAPK(unpackDir); err != nil {
		return nil, b.fail("unpack", err)
	}

	// 5. Inject frontend assets
	if err := b.injectFrontend(unpackDir, in.ProjectDir); err != nil {
		return nil, b.fail("inject-frontend", err)
	}

	// 6. Patch all strings in binary files (AXML, resources.arsc, etc.)
	if err := patchAllStrings(unpackDir, pkgName, in.AppName, verName, verCode); err != nil {
		return nil, b.fail("patch-strings", err)
	}

	// 7. Place icons
	if len(in.IconWebP) > 0 {
		if err := b.placeIcons(unpackDir, in.IconWebP); err != nil {
			return nil, b.fail("icons", err)
		}
	}

	// 8. Repack unsigned APK
	unsignedPath := filepath.Join(taskDir, "unsigned.apk")
	if err := b.repackAPK(unpackDir, unsignedPath); err != nil {
		return nil, b.fail("repack", err)
	}

	// 9. Sign APK
	signedName := fmt.Sprintf("%s_v%s.apk", sanitizeFilename(in.AppName), verName)
	signedPath := filepath.Join(b.OutputDir, signedName)
	if err := b.signAPK(unsignedPath, signedPath); err != nil {
		return nil, b.fail("sign", err)
	}

	// 10. Cleanup
	if !b.KeepWorkDir {
		os.RemoveAll(taskDir)
	}

	elapsed := time.Since(start)
	b.logf("=== Build complete in %v ===", elapsed)
	b.logf("Output: %s", signedPath)

	return &BuildResult{
		APKPath:     signedPath,
		PackageName: pkgName,
		VersionName: verName,
		VersionCode: verCode,
		Log:         b.logBuf.String(),
	}, nil
}

// ── Step implementations ──────────────────────────────────────────────────

func (b *Builder) validateInput(in BuildInput) error {
	if in.ProjectDir == "" {
		return fmt.Errorf("project directory is required")
	}
	info, err := os.Stat(in.ProjectDir)
	if err != nil {
		return fmt.Errorf("project directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory")
	}

	// Check for index.html
	indexPath := filepath.Join(in.ProjectDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Search recursively
		var found bool
		filepath.Walk(in.ProjectDir, func(path string, fi os.FileInfo, err error) error {
			if err != nil || found {
				return nil
			}
			if strings.EqualFold(fi.Name(), "index.html") {
				found = true
			}
			return nil
		})
		if !found {
			return fmt.Errorf("no index.html found in project directory")
		}
	}

	if in.AppName == "" {
		return fmt.Errorf("app name is required")
	}

	if _, err := os.Stat(b.TemplatePath); os.IsNotExist(err) {
		return fmt.Errorf("template APK not found: %s", b.TemplatePath)
	}
	return nil
}

func (b *Builder) createWorkspace() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	taskDir := filepath.Join(b.WorkDir, "task_"+timestamp)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return "", fmt.Errorf("create workspace: %w", err)
	}
	b.logf("Workspace: %s", taskDir)
	return taskDir, nil
}

func (b *Builder) unpackAPK(dest string) error {
	b.logf("Unpacking template APK...")
	r, err := zip.OpenReader(b.TemplatePath)
	if err != nil {
		return fmt.Errorf("open template apk: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open %s: %w", f.Name, err)
		}

		out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("create %s: %w", fpath, err)
		}

		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return fmt.Errorf("write %s: %w", f.Name, err)
		}
	}
	b.logf("Unpacked %d files", len(r.File))
	return nil
}

func (b *Builder) injectFrontend(unpackDir, projectDir string) error {
	frontendDist := filepath.Join(unpackDir, "assets", "frontend", "dist")
	rootAssets := filepath.Join(unpackDir, "assets")
	os.MkdirAll(frontendDist, 0755)

	copied1 := 0
	copied2 := 0

	filepath.Walk(projectDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(projectDir, path)
		if err != nil {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		// 1. Copy to assets/frontend/dist/ (exact relative path preserved)
		dest1 := filepath.Join(frontendDist, rel)
		os.MkdirAll(filepath.Dir(dest1), 0755)
		if err := os.WriteFile(dest1, data, fi.Mode()); err != nil {
			return fmt.Errorf("write %s: %w", dest1, err)
		}
		copied1++

		// 2. Also copy to assets/ root, but flatten the project's "assets/" subdirectory.
		//    Vite puts JS/CSS under dist/assets/. For root-level loading, these need
		//    to be at assets/ directly (not assets/assets/).
		dest2Rel := rel
		// If the project uses a standard Vite structure with "assets/" subdirectory,
		// strip that prefix when copying to root
		// (e.g., "assets/foo.js" → "foo.js" at assets root)
		assetsPrefix := "assets" + string(filepath.Separator)
		if strings.HasPrefix(rel, assetsPrefix) {
			dest2Rel = rel[len(assetsPrefix):]
		}

		dest2 := filepath.Join(rootAssets, dest2Rel)
		os.MkdirAll(filepath.Dir(dest2), 0755)
		if err := os.WriteFile(dest2, data, fi.Mode()); err != nil {
			return fmt.Errorf("write %s: %w", dest2, err)
		}
		copied2++
		return nil
	})

	b.logf("Injected frontend: %d files (dist), %d files (root)", copied1, copied2)
	return nil
}

func (b *Builder) placeIcons(unpackDir string, icons map[string][]byte) error {
	// icons map: key = mipmap path relative to res/, value = webp bytes
	for relPath, data := range icons {
		dest := filepath.Join(unpackDir, "res", relPath)
		os.MkdirAll(filepath.Dir(dest), 0755)
		if err := os.WriteFile(dest, data, 0644); err != nil {
			return fmt.Errorf("write icon %s: %w", relPath, err)
		}
	}
	b.logf("Placed %d icon files", len(icons))
	return nil
}

func (b *Builder) repackAPK(unpackDir, outputPath string) error {
	b.logf("Repacking APK...")

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	cw := &countingWriter{w: out}
	zw := zip.NewWriter(cw)
	// Note: no defer zw.Close() — we close it explicitly before reading the count

	// Walk the unpacked directory and add all files
	err = filepath.Walk(unpackDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(unpackDir, path)
		if err != nil {
			return err
		}
		zipName := filepath.ToSlash(rel)

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		// APK requirements:
		//   - resources.arsc MUST be stored uncompressed and 4-byte aligned
		//   - AndroidManifest.xml should be stored uncompressed
		//   - Other binary files (classes.dex, .so, images) should use Store
		//   - Text-based files (JS, CSS, HTML) can use Deflate
		isARSC := strings.EqualFold(zipName, "resources.arsc")
		isManifest := strings.EqualFold(zipName, "androidmanifest.xml")
		isCompressed := strings.HasSuffix(zipName, ".js") ||
			strings.HasSuffix(zipName, ".css") ||
			strings.HasSuffix(zipName, ".html") ||
			strings.HasSuffix(zipName, ".htm")

		method := zip.Deflate
		if isARSC || isManifest || !isCompressed {
			method = zip.Store
		}

		// For stored (uncompressed) files, pad extra field to ensure
		// the file data starts at a 4-byte aligned offset.
		var extra []byte
		if method == zip.Store {
			// currentPos = bytes written so far (before this header)
			currentPos := cw.written
			// Header size = 30 (fixed) + len(zipName) + extraLen
			// Data needs: (currentPos + 30 + len(zipName) + extraLen) % 4 == 0
			headerOverhead := 30 + len(zipName)
			padding := (4 - (int(currentPos+int64(headerOverhead)) % 4)) % 4
			if padding > 0 {
				extra = make([]byte, padding)
			}
		}

		w, err := zw.CreateHeader(&zip.FileHeader{
			Name:   zipName,
			Method: method,
			Extra:  extra,
		})
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", zipName, err)
		}

		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("write %s: %w", zipName, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Close the zip writer before checking size
	if err := zw.Close(); err != nil {
		return fmt.Errorf("close zip: %w", err)
	}

	b.logf("Repacked to %s (%d bytes)", outputPath, cw.written)
	return nil
}

// countingWriter wraps an io.Writer and counts bytes written.
type countingWriter struct {
	w       io.Writer
	written int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.written += int64(n)
	return n, err
}

func (b *Builder) signAPK(unsignedPath, signedPath string) error {
	b.logf("Signing APK...")

	// Ensure META-INF entries from unpacking are removed (they won't validate)
	// The signer will regenerate them.

	unsignedData, err := os.ReadFile(unsignedPath)
	if err != nil {
		return fmt.Errorf("read unsigned apk: %w", err)
	}

	// Load or generate signing keys
	keys, err := b.loadSigningKeys()
	if err != nil {
		return fmt.Errorf("load signing keys: %w", err)
	}

	// Parse and sign
	z, err := apksign.NewZip(unsignedData)
	if err != nil {
		return fmt.Errorf("parse apk: %w", err)
	}

	z, err = z.Sign(keys)
	if err != nil {
		return fmt.Errorf("sign apk: %w", err)
	}

	// Verify
	v1Err := z.VerifyV1()
	v2Err := z.VerifyV2()
	if v1Err != nil {
		b.logf("V1 verify: %v", v1Err)
	}
	if v2Err != nil {
		b.logf("V2 verify FAILED: %v", v2Err)
	} else {
		b.logf("V2 signature verified OK")
	}
	if z.IsAPK {
		b.logf("Signed APK: IsAPK=%v V1=%v V2=%v", z.IsAPK, z.IsV1Signed, z.IsV2Signed)
	}

	signedData := z.Bytes()
	if err := os.WriteFile(signedPath, signedData, 0644); err != nil {
		return fmt.Errorf("write signed apk: %w", err)
	}

	// Zipalign the signed APK: stored entries must be 4-byte aligned
	if err := ZipalignViaRepack(signedPath); err != nil {
		b.logf("WARNING: zipalign failed: %v (APK may not install on R+)", err)
	} else {
		b.logf("Zipaligned OK")
	}

	// Re-sign V2 since zipalign changed local file headers (invalidating V2 block)
	alignedData, err := os.ReadFile(signedPath)
	if err == nil {
		alignedZ, err := apksign.NewZip(alignedData)
		if err == nil {
			alignedZ, err = alignedZ.SignV2(keys)
			if err == nil {
				alignedData = alignedZ.Bytes()
				os.WriteFile(signedPath, alignedData, 0644)
				// Verify final v2
				if err := alignedZ.VerifyV2(); err != nil {
					b.logf("WARNING: final v2 verification: %v", err)
				} else {
					b.logf("Final V2 signature verified OK")
				}
			}
		}
	}

	b.logf("Signed APK -> %s (%d bytes)", signedPath, len(signedData))
	return nil
}

func (b *Builder) loadSigningKeys() ([]*android.SigningCert, error) {
	keyPath := filepath.Join(b.KeysDir, "default.key")
	certPath := filepath.Join(b.KeysDir, "default.crt")

	// Generate keys if they don't exist
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		b.logf("Generating new signing key pair...")
		os.MkdirAll(b.KeysDir, 0700)

		if err := generateKeyPair(keyPath, certPath); err != nil {
			return nil, fmt.Errorf("generate keys: %w", err)
		}
		b.logf("Generated key: %s", keyPath)
	}

	return []*android.SigningCert{
		{
			SigningKey: android.SigningKey{
				KeyPath: keyPath,
				Type:    android.RSA,
				Hash:    android.SHA256,
			},
			CertPath: certPath,
		},
	}, nil
}

// generateKeyPair creates an RSA 2048-bit key and self-signed cert using crypto/x509.
func generateKeyPair(keyPath, certPath string) error {
	// We use crypto/x509 + crypto/rsa to generate and save in PEM format.
	// The apksign library expects PKCS1 or PKCS8 private key + X.509 cert.
	importCmd := fmt.Sprintf(
		`openssl req -x509 -newkey rsa:2048 -keyout "%s" -out "%s" -days 10000 -nodes -subj "/CN=Foil"`,
		keyPath, certPath)
	// Try openssl first; fall back to Go-native generation
	if err := runCommand(importCmd); err != nil {
		return goGenerateKeyPair(keyPath, certPath)
	}
	return nil
}

func runCommand(cmd string) error {
	// Simple command runner
	c := parseCommand(cmd)
	return c.Run()
}

// parseCommand splits a command string into name and args.
func parseCommand(cmd string) *execCmd {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return &execCmd{name: "echo"}
	}
	return &execCmd{name: parts[0], args: parts[1:]}
}

type execCmd struct {
	name string
	args []string
}

func (c *execCmd) Run() error {
	cmd := createExec(c.name, c.args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s: %s", c.name, err, string(output))
	}
	return nil
}

// fail logs and wraps an error with step context.
func (b *Builder) fail(step string, err error) error {
	b.logf("FAILED at %s: %v", step, err)
	return fmt.Errorf("builder.%s: %w", step, err)
}

func sanitizeFilename(name string) string {
	r := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return r.Replace(name)
}
