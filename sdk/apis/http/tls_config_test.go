package apis

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func TestHttpClientFactory_UsesConfigDisableSSLVerification(t *testing.T) {
	t.Setenv("UNMESHED_AGENT_DISABLE_SSL_VERIFICATION", "false")

	config := configs.NewClientConfig()
	config.SetDisableSSLVerification(true)

	client := NewHttpClientFactory(config).Create()
	transport, ok := client.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
}

func TestHttpRequestFactory_UsesEnvDisableSSLVerificationFallback(t *testing.T) {
	t.Setenv("UNMESHED_AGENT_DISABLE_SSL_VERIFICATION", "true")

	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")

	factory := NewHttpRequestFactory(config)
	transport, ok := factory.client.Client.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
}

func TestHttpClientFactory_LoadsRootCAsFromConfiguredDirectory(t *testing.T) {
	certDir := t.TempDir()
	cert := writeTestCertificate(t, filepath.Join(certDir, "server.crt"))
	if err := os.WriteFile(filepath.Join(certDir, "ignore.txt"), []byte("not-a-cert"), 0644); err != nil {
		t.Fatalf("failed to write non-cert test file: %v", err)
	}

	config := configs.NewClientConfig()
	config.SetCACertDirectory(certDir)

	client := NewHttpClientFactory(config).Create()
	transport, ok := client.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.NotNil(t, transport.TLSClientConfig.RootCAs)
	_, err := cert.Verify(x509.VerifyOptions{Roots: transport.TLSClientConfig.RootCAs})
	assert.NoError(t, err)
}

func TestHttpRequestFactory_LoadsRootCAsFromConfiguredDirectory(t *testing.T) {
	certDir := t.TempDir()
	cert := writeTestCertificate(t, filepath.Join(certDir, "server.crt"))

	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	config.SetCACertDirectory(certDir)

	factory := NewHttpRequestFactory(config)
	transport, ok := factory.client.Client.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.NotNil(t, transport.TLSClientConfig.RootCAs)
	_, err := cert.Verify(x509.VerifyOptions{Roots: transport.TLSClientConfig.RootCAs})
	assert.NoError(t, err)
}

func TestHttpClientFactory_PanicsForInvalidCACertDirectory(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetCACertDirectory(filepath.Join(t.TempDir(), "missing"))

	assert.PanicsWithError(t,
		"failed to access CA certificate directory \""+*config.GetCACertDirectory()+"\": stat "+*config.GetCACertDirectory()+": no such file or directory",
		func() {
			NewHttpClientFactory(config).Create()
		},
	)
}

func TestHttpClientFactory_PanicsWhenCACertDirectoryContainsNoCRTFiles(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetCACertDirectory(t.TempDir())

	assert.PanicsWithError(t,
		"no .crt files found in CA certificate directory \""+*config.GetCACertDirectory()+"\"",
		func() {
			NewHttpClientFactory(config).Create()
		},
	)
}

func TestHttpClientFactory_LeavesSSLVerificationEnabledByDefault(t *testing.T) {
	t.Setenv("UNMESHED_AGENT_DISABLE_SSL_VERIFICATION", "false")

	config := configs.NewClientConfig()

	client := NewHttpClientFactory(config).Create()
	transport, ok := client.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.Nil(t, transport.TLSClientConfig)
}

func writeTestCertificate(t *testing.T, certPath string) *x509.Certificate {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "unmeshed-test-ca",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		t.Fatalf("failed to write certificate: %v", err)
	}

	return cert
}
