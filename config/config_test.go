package config

import (
	"path/filepath"
	"testing"

	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/outputs/transport"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	absPath, err := filepath.Abs("../tests/files/")

	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	config := &ConfigSettings{}

	err = cfgfile.Read(config, absPath+"/config.yml")
	assert.Nil(t, err)

	urls := config.Httpbeat.Urls
	assert.Equal(t, 3, len(urls))

	assert.Equal(t, "http://example.org/1", urls[0].Url)
	assert.Equal(t, "get", urls[0].Method)
	assert.Equal(t, "@every 1m", urls[0].Cron)
	assert.Equal(t, "body", urls[0].Body)
	assert.Equal(t, "foo1", urls[0].BasicAuth.Username)
	assert.Equal(t, "bar1", urls[0].BasicAuth.Password)
	assert.Equal(t, "http://proxy:3128", urls[0].ProxyUrl)
	assert.Equal(t, int64(120), *urls[0].Timeout)
	assert.Equal(t, 2, len(urls[0].Headers))
	assert.Equal(t, 2, len(urls[0].Fields))
	assert.Equal(t, "jolokia", urls[0].DocumentType)
	assert.Equal(t, 1, len(urls[0].SSL.CAs))
	assert.Equal(t, "/etc/pki/root/ca.pem", urls[0].SSL.CAs[0])
	assert.Equal(t, "/etc/pki/client/cert.pem", urls[0].SSL.Certificate.Certificate)
	assert.Equal(t, "/etc/pki/client/cert.key", urls[0].SSL.Certificate.Key)
	assert.Equal(t, transport.VerifyFull, urls[0].SSL.VerificationMode)
	assert.Equal(t, 1, len(urls[0].SSL.CipherSuites))
	assert.Equal(t, 0, len(urls[0].SSL.CurveTypes))
	assert.Equal(t, transport.TLSVersionSSL30, urls[0].SSL.Versions[0])
	assert.Equal(t, "unflatten", urls[0].JsonDotMode)

	assert.Equal(t, "http://example.org/2", urls[1].Url)
	assert.Equal(t, "post", urls[1].Method)
	assert.Equal(t, "@every 2m", urls[1].Cron)
	assert.Equal(t, "", urls[1].BasicAuth.Username)
	assert.Equal(t, "", urls[1].BasicAuth.Password)
	assert.Equal(t, "", urls[1].ProxyUrl)
	assert.Equal(t, 0, len(urls[1].Headers))
	assert.Equal(t, 0, len(urls[1].Fields))
	assert.Equal(t, "", urls[1].DocumentType)
	assert.Equal(t, "replace", urls[1].JsonDotMode)
	assert.Equal(t, "", urls[1].JsonDotModeCharacter)

	assert.Equal(t, "http://example.org/2", urls[2].Url)
	assert.Equal(t, "post", urls[2].Method)
	assert.Equal(t, "@every 2m", urls[2].Cron)
	assert.Equal(t, "", urls[2].BasicAuth.Username)
	assert.Equal(t, "", urls[2].BasicAuth.Password)
	assert.Equal(t, "", urls[2].ProxyUrl)
	assert.Equal(t, 0, len(urls[2].Headers))
	assert.Equal(t, 0, len(urls[2].Fields))
	assert.Equal(t, "", urls[2].DocumentType)
	assert.Equal(t, "replace", urls[2].JsonDotMode)
	assert.Equal(t, "-", urls[2].JsonDotModeCharacter)
}
