package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"lets_config/internal/dpapi"
)

// certInfo stores certificate credentials, encrypted with DPAPI on disk.
type certInfo struct {
	CertPath     string `json:"certPath"`
	CertPassword string `json:"certPassword"`
	CertAlias    string `json:"certAlias"`
	KeyPassword  string `json:"keyPassword"`
}

// certStoragePath returns the path to the encrypted cert info file.
func certStoragePath() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, "unieditdept", "foil", "cert.dat")
}

// saveCertInfo encrypts and writes cert credentials to disk.
func saveCertInfo(certPath, certPassword, certAlias, keyPassword string) error {
	info := certInfo{
		CertPath:     certPath,
		CertPassword: certPassword,
		CertAlias:    certAlias,
		KeyPassword:  keyPassword,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	encrypted, err := dpapi.Protect(data)
	if err != nil {
		return err
	}
	// Ensure storage directory exists
	storageDir := filepath.Dir(certStoragePath())
	if err := os.MkdirAll(storageDir, 0700); err != nil {
		return err
	}
	return os.WriteFile(certStoragePath(), encrypted, 0600)
}

// loadCertInfo reads and decrypts cert credentials from disk.
func loadCertInfo() (*certInfo, error) {
	encrypted, err := os.ReadFile(certStoragePath())
	if err != nil {
		if os.IsNotExist(err) {
			return &certInfo{}, nil
		}
		return nil, err
	}
	decrypted, err := dpapi.Unprotect(encrypted)
	if err != nil {
		return nil, err
	}
	var info certInfo
	if err := json.Unmarshal(decrypted, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// clearCertInfo removes the encrypted cert info file.
func clearCertInfo() error {
	if err := os.Remove(certStoragePath()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
