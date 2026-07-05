package builder

import (
	"fmt"
	"strings"
	"testing"
)

func TestIsValidVersionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", true},
		{"single digit", "1", true},
		{"major.minor", "1.0", true},
		{"three parts", "1.2.3", true},
		{"long version", "2025.12.31.01", true},
		{"leading zero", "0.1.0", true},
		{"letters", "1.a", false},
		{"hyphen", "1-0", false},
		{"space", "1. 0", false},
		{"unicode", "1.零", false},
		{"emoji", "1.😀", false},
		{"special chars", "1.*", false},
		{"shell chars", "1;rm -rf /", false},
		{"path traversal", "../etc", false},
		{"newline", "1\n0", false},
		{"carriage return", "1\r0", false},
		{"null byte", "1\x000", false},
		{"very long", strings.Repeat("1.", 1000), true}, // all chars are '1' or '.' so technically valid; no length limit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidVersionName(tt.input)
			if got != tt.want {
				t.Errorf("isValidVersionName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidPackageName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", true},
		{"normal", "com.example.app", true},
		{"two segments", "com.app", true},
		{"with underscore", "com.my_app.app", true},
		{"with digits", "com.example2.app3", true},
		{"single segment", "com", false},
		{"empty segment", "com..app", false},
		{"leading digit", "com.2app", false},
		{"uppercase", "com.Example.App", false},
		{"hyphen", "com.my-app.app", false}, // Android doesn't allow hyphens
		{"special chars", "com.app$pecial", false},
		{"space", "com.my app", false},
		{"chinese", "com.我的.app", false},
		{"emoji", "com.😀.app", false},
		{"newline", "com.app\n.hack", false},
		{"sql injection", "com.';DROP TABLE;--", false},
		{"path traversal", "com.../etc", false},
		{"very long", "com." + strings.Repeat("a", 500), true}, // all lowercase letters = valid; no length limit
		{"valid long segment", "com." + strings.Repeat("a", 100), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPackageName(tt.input)
			if got != tt.want {
				t.Errorf("isValidPackageName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGeneratePackageName(t *testing.T) {
	tests := []struct {
		appName string
	}{
		{"HelloWorld"},
		{"我的应用"},
		{"TestApp123"},
		{"a"},
		{""},
		{"A very long application name that exceeds normal limits"},
		{"!@#$%^&*()"},
		{"test-app"},
		{"test_app"},
		{strings.Repeat("x", 200)},
	}

	for _, tt := range tests {
		t.Run(tt.appName, func(t *testing.T) {
			pkg := GeneratePackageName(tt.appName, nil, 26)
			if pkg == "" {
				t.Error("GeneratePackageName returned empty string")
			}
			if len(pkg) != 26 {
				t.Errorf("expected package name length 26, got %d: %q", len(pkg), pkg)
			}
			if !strings.HasPrefix(pkg, "com.") {
				t.Errorf("package name must start with 'com.': %q", pkg)
			}
			// Verify it passes our own validation
			if !isValidPackageName(pkg) {
				t.Errorf("generated package name %q fails isValidPackageName", pkg)
			}
		})
	}
}

func TestGeneratePackageNameUniqueness(t *testing.T) {
	// Generate 50 package names from the same app name and verify they're unique
	const count = 50
	pkgs := make(map[string]bool, count)
	for i := 0; i < count; i++ {
		pkg := GeneratePackageName("MyApp", nil, 26)
		if pkgs[pkg] {
			t.Errorf("duplicate package name generated: %q", pkg)
		}
		pkgs[pkg] = true
	}
}

func TestGenerateVersionName(t *testing.T) {
	v := GenerateVersionName()
	if len(v) < 2 || len(v) > 3 {
		t.Errorf("unexpected version name length: %q (%d chars)", v, len(v))
	}
	if !isValidVersionName(v) {
		t.Errorf("generated version name %q fails validation", v)
	}
}

func TestGenerateVersionCode(t *testing.T) {
	vc := GenerateVersionCode()
	if vc <= 0 {
		t.Errorf("version code must be positive, got %d", vc)
	}
	// Version code should fit in int32
	if vc < 0 {
		t.Errorf("version code overflow: %d", vc)
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My App", "My_App"},
		{"test/app", "test_app"},
		{"test\\app", "test_app"},
		{"test:app", "test_app"},
		{"test*app", "test_app"},
		{"test?app", "test_app"},
		{"test\"app", "test_app"},
		{"test<app", "test_app"},
		{"test>app", "test_app"},
		{"test|app", "test_app"},
		{"normal", "normal"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Fuzz tests

func FuzzIsValidVersionName(f *testing.F) {
	f.Add("1.0")
	f.Add("")
	f.Add("1.2.3")
	f.Fuzz(func(t *testing.T, input string) {
		// Must never panic
		result := isValidVersionName(input)
		// Empty must be valid
		if input == "" && !result {
			t.Errorf("empty string must be valid")
		}
	})
}

func FuzzIsValidPackageName(f *testing.F) {
	f.Add("com.example.app")
	f.Add("")
	f.Add("a.b")
	f.Fuzz(func(t *testing.T, input string) {
		// Must never panic
		result := isValidPackageName(input)
		// Empty must be valid
		if input == "" && !result {
			t.Errorf("empty string must be valid")
		}
		// Check it handles any Unicode without panic
		_ = result
	})
}

func FuzzGeneratePackageName(f *testing.F) {
	f.Add("MyApp")
	f.Add("")
	f.Add(strings.Repeat("x", 500))
	f.Fuzz(func(t *testing.T, appName string) {
		// Must never panic, even with garbage input
		pkg := GeneratePackageName(appName, nil, 26)
		if pkg == "" {
			t.Errorf("GeneratePackageName should never return empty")
		}
	})
}

func FuzzSanitizeFilename(f *testing.F) {
	f.Add("test")
	f.Add("a/b")
	f.Fuzz(func(t *testing.T, input string) {
		// Must never panic
		result := sanitizeFilename(input)
		_ = result
	})
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"a & b", "a &amp; b"},
		{"a < b", "a &lt; b"},
		{"a > b", "a &gt; b"},
		{"it's", "it&apos;s"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&apos;xss&apos;)&lt;/script&gt;"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := escapeXML(tt.input)
			if got != tt.want {
				t.Errorf("escapeXML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Benchmark tests

func BenchmarkGeneratePackageName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneratePackageName("MyApp", nil, 26)
	}
}

func BenchmarkIsValidPackageName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidPackageName("com.example.myapp")
	}
}

func BenchmarkIsValidVersionName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidVersionName("1.2.3")
	}
}

func BenchmarkSanitizeFilename(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sanitizeFilename("My App/test:file")
	}
}

// Test edge cases for GeneratePackageName with different lengths
func TestGeneratePackageNameDifferentLengths(t *testing.T) {
	lengths := []int{10, 26, 50, 100}
	for _, l := range lengths {
		t.Run(fmt.Sprintf("len=%d", l), func(t *testing.T) {
			pkg := GeneratePackageName("Test", nil, l)
			if len(pkg) != l {
				t.Errorf("expected length %d, got %d: %q", l, len(pkg), pkg)
			}
			if !strings.HasPrefix(pkg, "com.") {
				t.Errorf("package must start with com.: %q", pkg)
			}
		})
	}
}
