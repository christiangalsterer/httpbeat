package main

import (
	"path/filepath"
	"testing"

	"github.com/elastic/libbeat/cfgfile"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	absPath, err := filepath.Abs("./tests/files/")

	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	config := &ConfigSettings{}

	err = cfgfile.Read(config, absPath+"/config.yml")
	assert.Nil(t, err)

	urls := config.Httpbeat.Urls
	assert.Equal(t, 2, len(urls))

	assert.Equal(t, "http://example.org/1", urls[0].Url)
	assert.Equal(t, "get", urls[0].Method)
	assert.Equal(t, int64(5), *urls[0].Period)
	assert.Equal(t, "body", urls[0].Body)
	assert.Equal(t, "foo1", urls[0].Username)
	assert.Equal(t, "bar1", urls[0].Password)
	assert.Equal(t, "proxyUser", urls[0].ProxyUsername)
	assert.Equal(t, "proxyPass", urls[0].ProxyPassword)
	assert.Equal(t, "proxy", urls[0].ProxyHost)
	assert.Equal(t, "3128", urls[0].ProxyPort)
	assert.Equal(t, int64(120), *urls[0].Timeout)
	assert.Equal(t, 2, len(urls[0].Headers))
	assert.Equal(t, 2, len(urls[0].Fields))
	assert.Equal(t, "jolokia", urls[0].DocumentType)

	assert.Equal(t, "http://example.org/2", urls[1].Url)
	assert.Equal(t, "post", urls[1].Method)
	assert.Equal(t, int64(10), *urls[1].Period)
	assert.Equal(t, "", urls[1].Username)
	assert.Equal(t, "", urls[1].Password)
	assert.Equal(t, "", urls[1].ProxyUsername)
	assert.Equal(t, "", urls[1].ProxyPassword)
	assert.Equal(t, "", urls[1].ProxyHost)
	assert.Equal(t, "", urls[1].ProxyPort)
	assert.Equal(t, 0, len(urls[1].Headers))
	assert.Equal(t, 0, len(urls[1].Fields))
	assert.Equal(t, "", urls[1].DocumentType)
}

