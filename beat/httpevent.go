package beat

import (
	"time"
	"github.com/elastic/beats/libbeat/common"
)

type HttpEvent struct {
	ReadTime     time.Time
	DocumentType string
	Fields       map[string]string
	Request    	 Request
	Response     Response
}

type Request struct {
	Url         string `json:"url"`
	Method      string `json:"method"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        string `json:"body,omitempty"`
}

type Response struct {
	StatusCode    int `json:"statusCode"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          string `json:"body,omitempty"`
}

func (h *HttpEvent) ToMapStr() common.MapStr {
	event := common.MapStr{
		"@timestamp": common.Time(h.ReadTime),
		"type":       h.DocumentType,
		"request":    h.Request,
		"response":   h.Response,
	}

	if h.Fields != nil {
		event["fields"] = h.Fields
	}

	return event
}