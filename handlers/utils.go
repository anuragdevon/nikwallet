package handlers

import (
	"bytes"
	"encoding/json"
)

func JsonEqual(json1, json2 []byte) bool {
	var json1Interface, json2interface interface{}
	if err := json.Unmarshal(json1, &json1Interface); err != nil {
		return false
	}
	if err := json.Unmarshal(json2, &json2interface); err != nil {
		return false
	}
	return bytes.Equal(json1, json2)
}
