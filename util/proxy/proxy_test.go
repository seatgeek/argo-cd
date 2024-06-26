package proxy

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProxyEnvIfAbsent(t *testing.T) {
	t.Run("Existing proxy env variables", func(t *testing.T) {
		proxy := "https://proxy:5000"
		cmd := exec.Command("test")
		cmd.Env = []string{`http_proxy="https_proxy=https://env-proxy:8888"`, "key=val"}
		got := UpsertEnv(cmd, proxy)
		assert.EqualValues(t, []string{"key=val", httpProxy(proxy), httpsProxy(proxy)}, got)
	})
	t.Run("proxy env variables not found", func(t *testing.T) {
		proxy := "http://proxy:5000"
		cmd := exec.Command("test")
		cmd.Env = []string{"key=val"}
		got := UpsertEnv(cmd, proxy)
		assert.EqualValues(t, []string{"key=val", httpProxy(proxy), httpsProxy(proxy)}, got)
	})
}

func TestGetCallBack(t *testing.T) {
	t.Run("custom proxy present", func(t *testing.T) {
		proxy := "http://proxy:8888"
		url, err := GetCallback(proxy)(nil)
		assert.NoError(t, err)
		assert.Equal(t, proxy, url.String())
	})
	t.Run("custom proxy absent", func(t *testing.T) {
		proxyEnv := "http://proxy:8888"
		t.Setenv("http_proxy", "http://proxy:8888")
		url, err := GetCallback("")(httptest.NewRequest(http.MethodGet, proxyEnv, nil))
		assert.NoError(t, err)
		assert.Equal(t, proxyEnv, url.String())
	})
}
