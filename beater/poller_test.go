package beater

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapsMerging(t *testing.T) {
	map1 := map[string]interface{}{"key1": "value1", "commonKey": map[string]interface{}{"subkey1": "subvalue1", "commonSubkey": map[string]interface{}{"one": 1}}}
	map2 := map[string]interface{}{"subkey2": "subvalue2", "commonSubkey": map[string]interface{}{"two": 2}}

	mergeMaps(map1, map2, "commonKey")

	assert.Equal(t, "value1", map1["key1"])
	assert.Equal(t, "subvalue1", map1["commonKey"].(map[string]interface{})["subkey1"])
	assert.Equal(t, "subvalue2", map1["commonKey"].(map[string]interface{})["subkey2"])
	assert.Equal(t, 1, map1["commonKey"].(map[string]interface{})["commonSubkey"].(map[string]interface{})["one"])
	assert.Equal(t, 2, map1["commonKey"].(map[string]interface{})["commonSubkey"].(map[string]interface{})["two"])
}

func TestJsonUnflatten(t *testing.T) {
	assert.Equal(t, "plain text", unflat("plain text"))

	dot1 := unflat(map[string]interface{}{"dot.notation": "value"}).(map[string]interface{})
	assert.Equal(t, "value", dot1["dot"].(map[string]interface{})["notation"])

	dot2 := unflat(map[string]interface{}{"dot.one": 1, "dot.two": 2}).(map[string]interface{})
	assert.Equal(t, 1, dot2["dot"].(map[string]interface{})["one"])
	assert.Equal(t, 2, dot2["dot"].(map[string]interface{})["two"])
}
