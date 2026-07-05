package builder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"lets_config/internal/dpapi"
)

const dpapiMarker = "FOIL_DPAPI_V1"

// goGenerateKeyPair creates an RSA 2048-bit key and self-signed certificate.
// The private key is encrypted with Windows DPAPI before saving to disk.
func goGenerateKeyPair(keyPath, certPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate rsa key: %w", err)
	}

	// Serialize the private key to PKCS1 DER
	keyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// Encrypt with DPAPI
	encrypted, err := dpapi.Protect(keyDER)
	if err != nil {
		return fmt.Errorf("dpapi encrypt key: %w", err)
	}

	// Write: marker + encrypted blob
	if err := os.WriteFile(keyPath, append([]byte(dpapiMarker), encrypted...), 0600); err != nil {
		return fmt.Errorf("write key file: %w", err)
	}

	// Create self-signed certificate
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:   "Foil Auto-Generated Key",
			Organization: []string{"Foil"},
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("create cert file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}); err != nil {
		return fmt.Errorf("encode certificate: %w", err)
	}

	return nil
}

// decryptPrivateKey reads a key file that may be DPAPI-encrypted (starts with
// FOIL_DPAPI_V1) or plain PEM (backward compatibility).
func decryptPrivateKey(keyPath string) ([]byte, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// Check for DPAPI marker
	if len(data) > len(dpapiMarker) && string(data[:len(dpapiMarker)]) == dpapiMarker {
		encrypted := data[len(dpapiMarker):]
		return dpapi.Unprotect(encrypted)
	}

	// Fallback: plain PEM — return raw data as-is
	return data, nil
}
