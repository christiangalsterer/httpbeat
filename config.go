package main

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
	ProxyHost string `yaml:"proxyHost"`
	ProxyPort string `yaml:"proxyPort"`
	ProxyUsername string `yaml:"proxyUsername"`
	ProxyPassword string `yaml:"proxyPassword"`
	Timeout *int64
}

type ConfigSettings struct {
	Httpbeat HttpbeatConfig
}
