//go:build windows

package main

import (
	"os"
	"path/filepath"
	"regexp"
	"unsafe"

	"golang.org/x/sys/windows"
)

// DesktopPath reads the real Desktop path from the Windows registry
// (HKCU\…\User Shell Folders\Desktop), falling back to %USERPROFILE%\Desktop.
func DesktopPath() string {
	var path windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_CURRENT_USER,
		windows.StringToUTF16Ptr(`Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`),
		0,
		windows.KEY_READ,
		&path,
	)
	if err != nil {
		return filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
	}
	defer windows.RegCloseKey(path)

	var buf [1024]uint16
	var typ uint32
	var n uint32 = 1024
	err = windows.RegQueryValueEx(
		path,
		windows.StringToUTF16Ptr("Desktop"),
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&buf[0])),
		&n,
	)
	if err != nil || (typ != windows.REG_SZ && typ != windows.REG_EXPAND_SZ) {
		return filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
	}
	desktop := windows.UTF16ToString(buf[:])
	// Windows registry uses %VAR% syntax; os.ExpandEnv only handles $VAR
	desktop = expandWindowsEnv(desktop)
	return desktop
}

// expandWindowsEnv replaces %VAR% patterns with environment variable values.
// os.ExpandEnv only handles $VAR syntax, but Windows registry uses %VAR%.
var envPattern = regexp.MustCompile(`%([^%]+)%`)

func expandWindowsEnv(s string) string {
	return envPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := match[1 : len(match)-1]
		if v := os.Getenv(name); v != "" {
			return v
		}
		return match
	})
}
