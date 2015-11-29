package main

import (
	"testing"
	"time"
	"github.com/elastic/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func TestFileEventToMapStr(t *testing.T) {
	// Test 'fields' is not present when it is nil.
	now := time.Now()
	event := HttpEvent{}
	event.DocumentType = "test"
	event.ReadTime = now
	mapStr := event.ToMapStr()
	_, found := mapStr["fields"]
	assert.False(t, found)
	assert.Equal(t, "test", mapStr["type"])
	assert.Equal(t, common.Time(now), mapStr["@timestamp"])
}
