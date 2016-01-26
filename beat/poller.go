package beat

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	"encoding/json"
)

type Poller struct {
	httpbeat      *Httpbeat
	config        UrlConfig
	cron          string
}

func NewPooler(httpbeat *Httpbeat, config UrlConfig) *Poller {
	poller := &Poller{
		httpbeat: httpbeat,
		config:   config,
	}

	return poller
}

func (p *Poller) Run() {

	// Setup DocumentType
	if p.config.DocumentType == "" {
		p.config.DocumentType = DefaultDocumentType
	}

	//init the cron schedule
	if p.config.Cron != "" {
		p.cron = p.config.Cron
	} else {
		p.cron = DefaultCron
	}

	cron := cron.New()
	cron.AddFunc(p.config.Cron, func() { p.runOneTime() })
	cron.Start()
}

func (p *Poller) runOneTime() error {
	request := gorequest.New()
	url := p.config.Url
	method := p.config.Method

	switch method {
	case "get":
		request.Get(url)
	case "delete":
		request.Delete(url)
	case "head":
		request.Head(url)
	case "patch":
		request.Patch(url)
	case "post":
		request.Post(url)
	case "put":
		request.Put(url)
	default:
		return fmt.Errorf("Unsupported HTTP method %g", method)
	}

	// set timeout
	if p.config.Timeout != nil {
		request.Timeout(time.Duration(*p.config.Timeout) * time.Second)
	} else {
		request.Timeout(DefaultTimeout)
	}

	// set authentication
	if p.config.BasicAuth.Username != "" && p.config.BasicAuth.Password != "" {
		request.BasicAuth.Username = p.config.BasicAuth.Username
		request.BasicAuth.Password = p.config.BasicAuth.Password
	}

	// set tls config
	useTLS := (p.config.TLS != nil)
	if useTLS {
		var err error
		var tlsConfig *tls.Config
		tlsConfig, err = outputs.LoadTLSConfig(p.config.TLS)
		if err != nil {
			return err
		}
		request.TLSClientConfig(tlsConfig)
	}

	// set body
	if p.config.Body !="" {
		switch method {
		case "patch", "post", "put":
			request.SendString(p.config.Body)
		default:
		}
	}

	// set headers
	request.Header = p.config.Headers

	// set proxy
	if p.config.ProxyUrl != "" {
		request.Proxy(p.config.ProxyUrl)
	}

	logp.Debug("Httpbeat", "Executing HTTP request: %v", request)
	now := time.Now()
	resp, body, errs:= request.End();

	if errs != nil {
		logp.Err("An error occured while executing HTTP request: %v", errs)
		return fmt.Errorf("An error occured while executing HTTP request: %v", errs)
	}

	requestEvent := Request{
		Url: url,
		Method: method,
		Headers: p.config.Headers,
		Body: p.config.Body,
	}

	var jsonBody map[string]interface{}
	if json.Unmarshal([]byte(body), &jsonBody) != nil {
		jsonBody = nil
	} else {
		jsonBody = unflat(jsonBody).(map[string]interface{})
	}

	responseEvent := Response{
		StatusCode:    resp.StatusCode,
		Headers:       p.GetResponseHeader(resp),
		Body:          body,
		JsonBody:      jsonBody,
	}

	event := HttpEvent{
		ReadTime:     now,
		DocumentType: p.config.DocumentType,
		Fields:       p.config.Fields,
		Request:      requestEvent,
		Response:     responseEvent,
	}

	p.httpbeat.events.PublishEvent(event.ToMapStr())

	return nil
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

func (p *Poller) Stop() {
}

func (p *Poller) GetResponseHeader(response gorequest.Response) map[string]string {

	responseHeader := make(map[string]string)
	for k,v := range response.Header {
		value := ""
		for _,h := range v {
			value+=h+" ,"
		}
		value = strings.TrimRight(value, " ,")
		responseHeader[k] = value
	}
	return responseHeader
}