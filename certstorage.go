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
	return os.WriteFile(certStoragePath(), encrypted, 0600)
}

// loadCertInfo reads and decrypts cert credentials from disk.
func loadCertInfo() (certPath, certPassword, certAlias, keyPassword string, err error) {
	encrypted, err := os.ReadFile(certStoragePath())
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", "", "", nil
		}
		return "", "", "", "", err
	}
	decrypted, err := dpapi.Unprotect(encrypted)
	if err != nil {
		return "", "", "", "", err
	}
	var info certInfo
	if err := json.Unmarshal(decrypted, &info); err != nil {
		return "", "", "", "", err
	}
	return info.CertPath, info.CertPassword, info.CertAlias, info.KeyPassword, nil
}

// clearCertInfo removes the encrypted cert info file.
func clearCertInfo() error {
	if err := os.Remove(certStoragePath()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
