package apis

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func TestHttpRequestFactory_UsesOpenShellTokenWithoutHashingAndLogsValue(t *testing.T) {
	openShellSupportLogOnce = sync.Once{}
	t.Cleanup(func() {
		openShellSupportLogOnce = sync.Once{}
	})

	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	config.SetAuthToken("openshell:resolve:unmeshed-api-key")

	var logOutput bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logOutput)
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
	})

	factory := NewHttpRequestFactory(config)
	NewHttpRequestFactory(config)

	assert.Equal(t, "Bearer openshell:resolve:unmeshed-api-key", factory.bearerValue)
	assert.Equal(t, 1, strings.Count(logOutput.String(), "OpenShell auth token resolution is supported using the openshell:resolve: prefix"))
	assert.Contains(t, logOutput.String(), "Auth token uses openshell:resolve; sending value without SDK hashing: openshell:resolve:unmeshed-api-key")
}
