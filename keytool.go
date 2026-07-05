package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// listKeystoreAliases uses the bundled keytool to list all aliases in a keystore.
func listKeystoreAliases(keystorePath, storePass string) ([]string, error) {
	// Resolve keytool path relative to the bundled JRE
	assetDir := AssetsDir()
	keytoolPath := filepath.Join(assetDir, "jre-minimal", "bin", "keytool.exe")

	// Try keytool.exe on Windows
	if _, err := os.Stat(keytoolPath); err != nil {
		keytoolPath = filepath.Join(assetDir, "jre-minimal", "bin", "keytool")
		if _, err := os.Stat(keytoolPath); err != nil {
			return nil, fmt.Errorf("keytool not found in bundled JRE")
		}
	}

	// Write password to temp file to avoid exposing it in process list
	passDir, err := os.MkdirTemp("", "foil-keytool-pass-*")
	if err != nil {
		return nil, fmt.Errorf("create password temp dir: %w", err)
	}
	defer os.RemoveAll(passDir)

	passFile := filepath.Join(passDir, "store_pass.txt")
	if err := os.WriteFile(passFile, []byte(storePass), 0600); err != nil {
		return nil, fmt.Errorf("write password file: %w", err)
	}

	cmd := exec.Command(keytoolPath, "-list", "-keystore", keystorePath,
		"-storepass:file", passFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("keytool failed: %w\n%s", err, out.String())
	}

	// Parse output: aliases listed as "Alias name: <name>" lines
	var aliases []string
	for _, line := range strings.Split(out.String(), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Alias name:") {
			alias := strings.TrimSpace(strings.TrimPrefix(line, "Alias name:"))
			if alias != "" {
				aliases = append(aliases, alias)
			}
		}
	}

	if len(aliases) == 0 {
		return nil, fmt.Errorf("no aliases found in keystore")
	}

	return aliases, nil
}
