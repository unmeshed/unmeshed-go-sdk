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
	openShellTokenLogOnce = sync.Once{}
	t.Cleanup(func() {
		openShellTokenLogOnce = sync.Once{}
	})

	config := configs.NewClientConfig()
	config.SetClientID("test-client")

	var logOutput bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logOutput)
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
	})

	config.SetAuthToken("ordinary-token")
	NewHttpRequestFactory(config)
	assert.NotContains(t, logOutput.String(), "openshell:resolve")

	config.SetAuthToken("openshell:resolve:unmeshed-api-key")
	factory := NewHttpRequestFactory(config)
	NewHttpRequestFactory(config)

	assert.Equal(t, "Bearer openshell:resolve:unmeshed-api-key", factory.bearerValue)
	assert.Equal(t, 1, strings.Count(logOutput.String(), "Auth token uses openshell:resolve; sending value without SDK hashing: openshell:resolve:unmeshed-api-key"))
}
