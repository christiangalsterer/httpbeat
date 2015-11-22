package main

import (
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/logp"
	"github.com/elastic/libbeat/publisher"
    "github.com/parnurzeal/gorequest"
	"fmt"
)

type Httpbeat struct {
	done                 chan struct{}
	period               time.Duration
	HbConfig             ConfigSettings
	events               publisher.Client
}

func (h *Httpbeat) Config(b *beat.Beat) error {

	err := cfgfile.Read(&h.HbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	//init the period
	if h.HbConfig.Httpbeat.Period != nil {
		h.period = time.Duration(*h.HbConfig.Httpbeat.Period) * time.Second
	} else {
		h.period = 1 * time.Second
	}

	logp.Info("httpbeat", "Init httpbeat")

	return nil
}

func (h *Httpbeat) Setup(b *beat.Beat) error {
	h.events = b.Events

	return nil
}

func (h *Httpbeat) Run(b *beat.Beat) error {
	var err error

	ticker := time.NewTicker(h.period)
	defer ticker.Stop()

	//main loop
	for {
		select {
		case <-h.done:
			return nil
		case <-ticker.C:
		}

		timerStart := time.Now()
		h.runOneTime(b)
		timerEnd := time.Now()

		duration := timerEnd.Sub(timerStart)
		if duration.Nanoseconds() > h.period.Nanoseconds() {
			logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
		}
	}

	return err
}

func (h *Httpbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (h *Httpbeat) Stop() {
	close(h.done)
}

func (h *Httpbeat) runOneTime(b *beat.Beat) error {
	request := gorequest.New()
	url := h.HbConfig.Httpbeat.Url
	method := h.HbConfig.Httpbeat.Method

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
	if h.HbConfig.Httpbeat.Timeout != nil {
		request.Timeout(time.Duration(*h.HbConfig.Httpbeat.Timeout) * time.Second)
	} else {
		request.Timeout(60 * time.Second)
	}

	// set authentication
	if h.HbConfig.Httpbeat.Username != "" && h.HbConfig.Httpbeat.Password != "" {
		request.BasicAuth.Username = h.HbConfig.Httpbeat.Username
		request.BasicAuth.Password = h.HbConfig.Httpbeat.Password
	}

	// set body
	if h.HbConfig.Httpbeat.Body !="" {
		request.SendString(h.HbConfig.Httpbeat.Body)
	}

	// set headers
	request.Header = h.HbConfig.Httpbeat.Headers

	// set proxy
	proxyUrl := h.GetProxyUrl()
	if proxyUrl != "" {
		request.Proxy(proxyUrl)
	}

	logp.Debug("Httpbeat", "Trying to make the following HTTP request: %v", request)
	resp, body, errs:= request.End();

	if errs != nil {
		logp.Err("An error occured while executing HTTP request: %v", errs)
		return fmt.Errorf("An error occured while executing HTTP request: %v", errs)
	}

	requestEvent := Request{
		Url: url,
		Headers: h.HbConfig.Httpbeat.Headers,
		Body: h.HbConfig.Httpbeat.Body,
	}

	responseEvent := Response{
		StatusCode:    resp.StatusCode,
		//Headers:       h.GetResponseHeader(resp),
		ContentLength: resp.ContentLength,
		Body:          body,
	}

	event := HttpEvent{
		DocumentType: "httpbeat",
		Request:      requestEvent,
		Response:     responseEvent,
	}

	h.events.PublishEvent(event.ToMapStr())

	return nil
}

func (h *Httpbeat) GetProxyUrl() (string) {
	proxyUrl := ""
	if h.HbConfig.Httpbeat.ProxyHost != "" && h.HbConfig.Httpbeat.ProxyPort != "" {
		proxyUrl = h.HbConfig.Httpbeat.ProxyHost + h.HbConfig.Httpbeat.ProxyPort
	}

	return proxyUrl;
}

/*
func (h *Httpbeat) GetResponseHeader(response gorequest.Response) map[string]string {
	responseHeader := make(map[string]string)
	for k,v := range response.Header {
		responseHeader[k] = v
	}
	return responseHeader
}
*/