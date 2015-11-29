package main

import (
	"testing"
	"time"
	"github.com/elastic/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func TestHttpEventToMapStr(t *testing.T) {
	now := time.Now()
	fields := make(map[string]string)
	fields["field1"] = "value1"
	fields["field2"] = "value2"
	request := Request{}
	request.Url = "www.example.org"
	headers := make(map[string]string)
	headers["header1"] = "value1"
	request.Headers = headers
	request.Body = "Body"
	request.Method = "get"

	event := HttpEvent{}
	event.Fields = fields
	event.DocumentType = "test"
	event.ReadTime = now
	event.Request = request
	mapStr := event.ToMapStr()
	_, fieldsExist := mapStr["fields"]
	assert.True(t, fieldsExist)
	_, requestExist := mapStr["request"]
	assert.True(t, requestExist)
	assert.Equal(t, "test", mapStr["type"])
	assert.Equal(t, common.Time(now), mapStr["@timestamp"])
}

func TestHttpEventToMapStrWIthEmptyFields(t *testing.T) {
	event := HttpEvent{}
	mapStr := event.ToMapStr()
	_, fieldsExist := mapStr["fields"]
	assert.False(t, fieldsExist)
}
