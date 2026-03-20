package apis

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func sslVerificationDisabled(clientConfig *configs.ClientConfig) bool {
	if clientConfig != nil && clientConfig.IsDisableSSLVerification() {
		return true
	}

	return strings.EqualFold(strings.TrimSpace(os.Getenv("UNMESHED_AGENT_DISABLE_SSL_VERIFICATION")), "true")
}

func buildTLSConfig(clientConfig *configs.ClientConfig) *tls.Config {
	if clientConfig == nil {
		return nil
	}

	tlsConfig := &tls.Config{}
	configured := false

	if sslVerificationDisabled(clientConfig) {
		tlsConfig.InsecureSkipVerify = true
		configured = true
	}

	caCertDirectory := clientConfig.GetCACertDirectory()
	if caCertDirectory != nil && strings.TrimSpace(*caCertDirectory) != "" {
		rootCAs, err := loadRootCAsFromDirectory(strings.TrimSpace(*caCertDirectory))
		if err != nil {
			log.Printf("INFO: skipping custom CA certificate directory: %v", err)
		} else {
			tlsConfig.RootCAs = rootCAs
			configured = true
		}
	}

	if !configured {
		return nil
	}

	return tlsConfig
}

func loadRootCAsFromDirectory(directoryPath string) (*x509.CertPool, error) {
	info, err := os.Stat(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access CA certificate directory %q: %w", directoryPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("CA certificate path %q must be a directory", directoryPath)
	}

	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate directory %q: %w", directoryPath, err)
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	loadedCerts := 0
	for _, entry := range entries {
		extension := strings.ToLower(filepath.Ext(entry.Name()))
		if entry.IsDir() || (extension != ".crt" && extension != ".pem") {
			continue
		}

		certPath := filepath.Join(directoryPath, entry.Name())
		certPEM, err := os.ReadFile(certPath)
		if err != nil {
			log.Printf("WARN: failed to read CA certificate %q: %v", certPath, err)
			continue
		}

		if ok := rootCAs.AppendCertsFromPEM(certPEM); !ok {
			log.Printf("WARN: failed to append CA certificate %q", certPath)
			continue
		}
		loadedCerts++
	}

	if loadedCerts == 0 {
		return nil, fmt.Errorf("no .crt or .pem files found in CA certificate directory %q", directoryPath)
	}

	return rootCAs, nil
}
