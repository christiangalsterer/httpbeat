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
	request       *gorequest.SuperAgent
}

func NewPooler(httpbeat *Httpbeat, config UrlConfig) *Poller {
	poller := &Poller{
		httpbeat: httpbeat,
		config:   config,
		//request:  gorequest.New(),
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
	if  p.request == nil {
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

	// set timeout
	if p.config.Timeout != nil {
		p.request.Timeout(time.Duration(*p.config.Timeout) * time.Second)
	} else {
		p.request.Timeout(DefaultTimeout)
	}

	// set authentication
	if p.config.BasicAuth.Username != "" && p.config.BasicAuth.Password != "" {
		p.request.BasicAuth.Username = p.config.BasicAuth.Username
		p.request.BasicAuth.Password = p.config.BasicAuth.Password
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
	resp, body, errs := p.request.End();

	if errs != nil {
		p.request = nil
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