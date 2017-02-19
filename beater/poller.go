package beater

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/christiangalsterer/httpbeat/config"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/transport"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	"strings"
	"time"
)

type Poller struct {
	httpbeat *Httpbeat
	config   config.HostConfig
	schedule string
	documentType string
	jsonDotModeCharacter string
	outputFormat string
	timeout string
	request  *gorequest.SuperAgent
}

func NewPooler(httpbeat *Httpbeat, config config.HostConfig) *Poller {
	poller := &Poller{
		httpbeat: httpbeat,
		config:   config,
	}

	return poller
}

func (p *Poller) Run() {

	// setup default config
	p.schedule = config.DefaultSchedule
	p.documentType = config.DefaultDocumentType
	p.jsonDotModeCharacter = config.DefaultJsonDotModeCharacter
	p.outputFormat = config.DefaultOutputFormat
	p.timeout = config.DefaultTimeout

	// setup DocumentType
	if p.config.DocumentType != "" {
		p.documentType = p.config.DocumentType
	}

	// setuo cron schedule
	if p.config.Schedule != "" {
		logp.Debug("Httpbeat", "Use schedule: [%w]", p.config.Schedule)
		p.schedule = p.config.Schedule
	}

	// setup jsonDotModeCharacter
	if p.config.JsonDotModeCharacter != "" {
		p.jsonDotModeCharacter = p.config.JsonDotModeCharacter

	}

	// setup timeout
	if p.config.Timeout != nil {
		p.request.Timeout(time.Duration(*p.config.Timeout) * time.Second)
	}

	cron := cron.New()
	cron.AddFunc(p.schedule, func() { p.runOneTime() })
	cron.Start()
}

func (p *Poller) runOneTime() error {
	if p.request == nil {
		p.request = gorequest.New()
	}

	url := p.config.Url
	method := p.config.Method
	switch method {
	case "get":
		p.request.Get(url)
	case "delete":
		p.request.Delete(url)
	case "head":
		p.request.Head(url)
	case "patch":
		p.request.Patch(url)
	case "post":
		p.request.Post(url)
	case "put":
		p.request.Put(url)
	default:
		return fmt.Errorf("Unsupported HTTP method %g", method)
	}

	outputFormat := p.config.OutputFormat
	switch outputFormat {
	case "":
		outputFormat = config.DefaultOutputFormat
	case "string":
	case "json":
		break
	default:
		return fmt.Errorf("Unsupported output format %g", outputFormat)
	}

	// set authentication
	if p.config.BasicAuth.Username != "" && p.config.BasicAuth.Password != "" {
		p.request.BasicAuth.Username = p.config.BasicAuth.Username
		p.request.BasicAuth.Password = p.config.BasicAuth.Password
	}

	// set tls config
	useTLS := (p.config.SSL != nil)
	if useTLS {
		var err error
		var tlsConfig *tls.Config
		var tlsC *transport.TLSConfig
		//tlsConfig, err = outputs.LoadTLSConfig(p.config.TLS)
		tlsC, err = outputs.LoadTLSConfig(p.config.SSL)
		if err != nil {
			return err
		}
		tlsConfig = convertTLSConfig(tlsC)
		p.request.TLSClientConfig(tlsConfig)
	}

	// set body
	if p.config.Body != "" {
		switch method {
		case "patch", "post", "put":
			p.request.SendString(p.config.Body)
		default:
		}
	}

	// set headers
	p.request.Header = p.config.Headers

	// set proxy
	if p.config.ProxyUrl != "" {
		p.request.Proxy(p.config.ProxyUrl)
	}

	logp.Debug("Httpbeat", "Executing HTTP request: %v", p.request)
	now := time.Now()
	resp, body, errs := p.request.End()

	if errs != nil {
		p.request = nil
		logp.Err("An error occurred while executing HTTP request: %v", errs)
		return fmt.Errorf("An error occurred while executing HTTP request: %v", errs)
	}

	requestEvent := Request{
		Url:     url,
		Method:  method,
		Headers: p.config.Headers,
		Body:    p.config.Body,
	}

	var jsonBody map[string]interface{}

	responseEvent := Response{
		StatusCode: resp.StatusCode,
		Headers:    p.GetResponseHeader(resp),
	}

	if outputFormat == "string" {
		responseEvent.Body = body
	} else {
		if outputFormat == "json" {
			decoder := json.NewDecoder(strings.NewReader(body))
			decoder.UseNumber()
			errs := decoder.Decode(&jsonBody)
			if errs != nil {
				jsonBody = nil
				logp.Err("An error occurred while marshalling response to JSON: %w", errs)
			} else {
				if p.config.JsonDotMode == "unflatten" {
					jsonBody = unflat(jsonBody).(map[string]interface{})
				} else if p.config.JsonDotMode == "replace" {
					jsonBody = replaceDots(jsonBody, p.jsonDotModeCharacter).(map[string]interface{})
				}
			}
			responseEvent.JsonBody = jsonBody
		}
	}

	event := HttpEvent{
		ReadTime:     now,
		DocumentType: p.documentType,
		Fields:       p.config.Fields,
		Request:      requestEvent,
		Response:     responseEvent,
	}

	p.httpbeat.client.PublishEvent(event.ToMapStr())

	return nil
}

func replaceDots(data interface{}, jsonDotModeCharacter string) interface{} {
	switch data.(type) {
	case map[string]interface{}:
		result := map[string]interface{}{}
		for key, value := range data.(map[string]interface{}) {
			result[strings.Replace(key, ".", jsonDotModeCharacter, -1)] = replaceDots(value, jsonDotModeCharacter)
		}
		return result
	default:
		return data
	}
}

func unflat(data interface{}) interface{} {
	switch data.(type) {
	case map[string]interface{}:
		result := map[string]interface{}{}
		for key, value := range data.(map[string]interface{}) {
			parts := strings.SplitN(key, ".", 2)
			if len(parts) < 2 {
				result[key] = unflat(value)
				continue
			} else {
				mergeMaps(result, unflat(map[string]interface{}{parts[1]: value}).(map[string]interface{}), parts[0])
			}
		}
		return result
	default:
		return data
	}
}

func mergeMaps(first map[string]interface{}, second map[string]interface{}, key string) {
	existingValue, exists := first[key]
	if !exists {
		first[key] = second
	} else {
		switch existingValue.(type) {
		case map[string]interface{}:
			for k, v := range second {
				switch v.(type) {
				case map[string]interface{}:
					mergeMaps(existingValue.(map[string]interface{}), v.(map[string]interface{}), k)
				default:
					existingValue.(map[string]interface{})[k] = v
				}
			}
		default:
			first[key] = second
		}
	}
}

func convertTLSConfig(config *transport.TLSConfig) *tls.Config {
	return &tls.Config{
		Certificates:       config.Certificates,
		CipherSuites:       config.CipherSuites,
		RootCAs:            config.RootCAs,
		CurvePreferences:   config.CurvePreferences,
		InsecureSkipVerify: config.Verification != transport.VerifyFull,
	}
}

func (p *Poller) Stop() {
}

func (p *Poller) GetResponseHeader(response gorequest.Response) map[string]string {

	responseHeader := make(map[string]string)
	for k, v := range response.Header {
		value := ""
		for _, h := range v {
			value += h + " ,"
		}
		value = strings.TrimRight(value, " ,")
		responseHeader[k] = value
	}
	return responseHeader
}
