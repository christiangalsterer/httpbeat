package beat

import (
	"time"
	"github.com/elastic/beats/libbeat/outputs"
)

// Defaults for config variables which are not set
const (
	DefaultPeriod              time.Duration = 1 * time.Second
	DefaultTimeout             time.Duration = 60 * time.Second
	DefaultDocumentType                      = "httpbeat"
)

type HttpbeatConfig struct {
	Urls []UrlConfig
}

type UrlConfig struct {
	Period *int64
	Url string
	Username string
	Password string
	Method string
	Body string
	Headers map[string]string
	ProxyHost string `yaml:"proxy_host"`
	ProxyPort string `yaml:"proxy_port"`
	ProxyUsername string `yaml:"proxy_username"`
	ProxyPassword string `yaml:"proxy_password"`
	Timeout *int64
	DocumentType string `yaml:"document_type"`
	Fields map[string]string `yaml:"fields"`
	TLS *outputs.TLSConfig
}

type ConfigSettings struct {
	Httpbeat HttpbeatConfig
}
