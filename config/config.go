package config

import (
	"github.com/elastic/beats/libbeat/outputs"
	"time"
)

// Defaults for config variables which are not set
const (
	DefaultCron         string        = "@every 1m"
	DefaultTimeout      time.Duration = 60 * time.Second
	DefaultDocumentType string        = "httpbeat"
)

type HttpbeatConfig struct {
	Urls []UrlConfig
}

type UrlConfig struct {
	Cron         string
	Url          string
	BasicAuth    BasicAuthenticationConfig `config:"basic_auth"`
	Method       string
	Body         string
	Headers      map[string]string
	ProxyUrl     string `config:"proxy_url"`
	Timeout      *int64
	DocumentType string            `config:"document_type"`
	Fields       map[string]string `config:"fields"`
	TLS          *outputs.TLSConfig
	JsonDotMode  string `config:"json_dot_mode"`
}

type BasicAuthenticationConfig struct {
	Username string
	Password string
}

type ConfigSettings struct {
	Httpbeat HttpbeatConfig
}
