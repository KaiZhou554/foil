package builder

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// ZipalignViaRepack ensures all stored entries in an APK are 4-byte aligned
// by re-reading the APK and creating a new ZIP with proper alignment.
// This invalidates the v2 signing block, so call SignV2 after this.
func ZipalignViaRepack(path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer r.Close()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Track byte position for alignment calculation.
	// We compute it manually: each entry's local header takes
	// 30 + len(name) + len(extra) bytes, followed by the file data.
	var pos int64

	for _, f := range r.File {
		if f.Method != 0 {
			// Compressed entry — copy as-is (no alignment needed)
			if err := zw.Copy(f); err != nil {
				return fmt.Errorf("copy %s: %w", f.Name, err)
			}
			pos += int64(30 + len(f.Name) + len(f.Extra) + int(f.CompressedSize64))
			continue
		}

		// Stored (uncompressed) entry — need 4-byte aligned data offset
		headerOverhead := int64(30 + len(f.Name))
		padding := (4 - ((pos + headerOverhead) % 4)) % 4
		var extra []byte
		if padding > 0 {
			extra = make([]byte, padding)
		}

		wh, err := zw.CreateHeader(&zip.FileHeader{
			Name:   f.Name,
			Method: zip.Store,
			Extra:  extra,
		})
		if err != nil {
			return fmt.Errorf("create %s: %w", f.Name, err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open %s: %w", f.Name, err)
		}
		n, err := io.Copy(wh, rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("write %s: %w", f.Name, err)
		}

		pos += int64(30 + len(f.Name) + len(extra) + int(n))
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("close zip: %w", err)
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}
