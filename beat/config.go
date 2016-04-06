package beat

import (
	"time"
	"github.com/elastic/beats/libbeat/outputs"
)

// Defaults for config variables which are not set
const (
	DefaultCron	string = "@every 1m"
	DefaultTimeout time.Duration = 60 * time.Second
	DefaultDocumentType string = "httpbeat"
)

type HttpbeatConfig struct {
	Urls []UrlConfig
}

type UrlConfig struct {
	Cron string
	Url string
	BasicAuth BasicAuthenticationConfig `yaml:"basic_auth"`
	Method string
	Body string
	Headers map[string]string
	ProxyUrl string `yaml:"proxy_url"`
	Timeout *int64
	DocumentType string `yaml:"document_type"`
	Fields map[string]string `yaml:"fields"`
	TLS *outputs.TLSConfig
	JsonDotMode string `yaml:"json_dot_mode"`
}

type BasicAuthenticationConfig struct {
	Username string
	Password string
}

type ConfigSettings struct {
	Httpbeat HttpbeatConfig
}
