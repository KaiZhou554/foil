//go:build windows

package dpapi

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modcrypt32  = windows.NewLazySystemDLL("crypt32.dll")
	procCryptProtectData   = modcrypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = modcrypt32.NewProc("CryptUnprotectData")
)

// dataBlob represents the CRYPTOAPI_BLOB structure used by DPAPI.
type dataBlob struct {
	cbData uint32
	pbData *byte
}

// Protect encrypts data using Windows DPAPI (user-specific, machine-local).
// The encrypted output is a binary blob that can only be decrypted by the
// same Windows user on the same machine.
func Protect(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, nil
	}

	in := &dataBlob{
		cbData: uint32(len(plaintext)),
		pbData: &plaintext[0],
	}
	var out dataBlob

	// CRYPTPROTECT_UI_FORBIDDEN = 0x1 — no UI prompt
	ret, _, err := procCryptProtectData.Call(
		uintptr(unsafe.Pointer(in)),
		0,                    // szDataDescr (optional)
		0,                    // optional entropy
		0,                    // reserved
		0,                    // prompt struct (nil = no UI)
		0x1,                  // dwFlags: CRYPTPROTECT_UI_FORBIDDEN
		uintptr(unsafe.Pointer(&out)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("CryptProtectData failed: %w", err)
	}

	result := make([]byte, out.cbData)
	copy(result, unsafe.Slice(out.pbData, out.cbData))

	// Free the allocated memory
	windows.LocalFree(windows.Handle(unsafe.Pointer(out.pbData)))

	return result, nil
}

// Unprotect decrypts data previously encrypted by Protect.
func Unprotect(encrypted []byte) ([]byte, error) {
	if len(encrypted) == 0 {
		return nil, nil
	}

	in := &dataBlob{
		cbData: uint32(len(encrypted)),
		pbData: &encrypted[0],
	}
	var out dataBlob

	ret, _, err := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(in)),
		0,                    // ppszDataDescr (optional)
		0,                    // optional entropy
		0,                    // reserved
		0,                    // prompt struct
		0x1,                  // dwFlags: CRYPTPROTECT_UI_FORBIDDEN
		uintptr(unsafe.Pointer(&out)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("CryptUnprotectData failed: %w", err)
	}

	result := make([]byte, out.cbData)
	copy(result, unsafe.Slice(out.pbData, out.cbData))

	windows.LocalFree(windows.Handle(unsafe.Pointer(out.pbData)))

	return result, nil
}
