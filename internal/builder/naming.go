package builder

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GeneratePackageName creates a unique Android package name from an app name.
// The output targets oldPkgLen chars to allow in-place AXML string replacement.
func GeneratePackageName(appName string, existingPkgNames []string, oldPkgLen int) string {
	sanitized := sanitizePackageSegment(appName)
	if sanitized == "" {
		sanitized = "app"
	}

	base := "com." + sanitized

	// Pad to exact length with hex chars
	if oldPkgLen > 0 {
		// Generate a unique suffix long enough to hit the target length
		needed := oldPkgLen - len(base)
		if needed > 0 {
			// Use timestamp-based filler
			filler := fmt.Sprintf("%x%x", time.Now().UnixNano(), time.Now().Unix())
			for len(base) < oldPkgLen {
				base += filler
				if len(base) > oldPkgLen {
					base = base[:oldPkgLen]
					break
				}
			}
		} else if needed < 0 {
			base = base[:oldPkgLen]
		}
	}

	// Ensure uniqueness
	for !isUnique(base, existingPkgNames) {
		time.Sleep(time.Nanosecond) // ensure different timestamp
		base = "com." + sanitized
		needed := oldPkgLen - len(base)
		if needed > 0 {
			filler := fmt.Sprintf("%x%x", time.Now().UnixNano(), time.Now().Unix())
			base += filler[:min(needed, len(filler))]
		}
		if len(base) > oldPkgLen {
			base = base[:oldPkgLen]
		}
	}

	return base
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sanitizePackageSegment(name string) string {
	// For non-Latin names (Chinese, etc.), use a hash-based fallback
	hasLatin := false
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			hasLatin = true
			break
		}
	}

	if !hasLatin {
		// Use a consistent hash of the name
		h := 0
		for _, r := range name {
			h = h*31 + int(r)
		}
		return fmt.Sprintf("app%x", h&0xFFFFF)
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	s := re.ReplaceAllString(strings.ToLower(name), "_")
	s = strings.Trim(s, "_")
	s = regexp.MustCompile(`^[0-9]+`).ReplaceAllString(s, "")
	if len(s) > 15 {
		s = s[:15]
	}
	return s
}

func isUnique(pkg string, existing []string) bool {
	for _, p := range existing {
		if p == pkg {
			return false
		}
	}
	return true
}

// GenerateVersionName returns a short version string (fits in "1.0"'s 3-char slot).
// Format: YDD where Y = last digit of year, DD = day of month.
// Example: "624" for June 24, 2026.
func GenerateVersionName() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d", now.Year()%10, now.Day())
}

// GenerateVersionCode returns an ever-increasing version code within int32 range.
// Format: last 9 digits of YYYYMMDDHHMM (safe for int32).
func GenerateVersionCode() int32 {
	now := time.Now()
	s := fmt.Sprintf("%d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute())
	// Take modulo to stay within int32 range
	var val int64
	fmt.Sscanf(s, "%d", &val)
	return int32(val % 1000000000)
}

// isValidVersionName checks that a version name contains only digits and dots
// (e.g. "1", "1.0", "2.3.1"). Empty string is allowed (triggers auto-generation).
func isValidVersionName(v string) bool {
	if v == "" {
		return true
	}
	for _, r := range v {
		if r != '.' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}
