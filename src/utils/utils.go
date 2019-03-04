package utils

import (
	"encoding/json"

	"github.com/golang/glog"
)

// Serialize ...
func Serialize(input interface{}) string {
	serial, err := json.Marshal(input)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return string(serial)
}

// Deserialize ...
func Deserialize(input string) interface{} {
	output := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &output)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return output
}
