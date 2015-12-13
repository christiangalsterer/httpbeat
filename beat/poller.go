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

	responseEvent := Response{
		StatusCode:    resp.StatusCode,
		Headers:       p.GetResponseHeader(resp),
		Body:          body,
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