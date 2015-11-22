package main

import (
	"time"
	"github.com/elastic/libbeat/common"
)

type HttpEvent struct {
	ReadTime     time.Time
	DocumentType string
	Fields       *map[string]string
	Request    	 Request
	Response     Response
}

type Request struct {
	Url         string `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        string `json:"body,omitempty"`
}

type Response struct {
	StatusCode    int `json:"statusCode,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string `json:"body,omitempty"`
	ContentLength int64 `json:"contentLength,omitempty"`
}

func (h *HttpEvent) ToMapStr() common.MapStr {
	event := common.MapStr{
		"@timestamp": common.Time(h.ReadTime),
		"type":       h.DocumentType,
		"request":    h.Request,
		"response":   h.Response,
	}

	return event
}