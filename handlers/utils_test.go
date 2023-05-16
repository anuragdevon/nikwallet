package handlers

import (
	"testing"
)

func TestJsonEqualToReturnTrueForEqualJson(t *testing.T) {
	jsonStr1 := `{"name":"John","age":30,"city":"New York"}`
	jsonStr2 := `{"name":"John","age":30,"city":"New York"}`

	json1 := []byte(jsonStr1)
	json2 := []byte(jsonStr2)

	if !JsonEqual(json1, json2) {
		t.Errorf("JsonEqual(%s, %s) = false, expected true", jsonStr1, jsonStr2)
	}
}
