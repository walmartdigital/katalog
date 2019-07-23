package utils

import (
	"encoding/json"
	"reflect"

	"k8s.io/klog"
)



// Serialize ...
func Serialize(input interface{}) string {
	serial, err := json.Marshal(input)
	if err != nil {
		klog.Error(err)
		return ""
	}
	return string(serial)
}

// Deserialize ...
func Deserialize(input string) interface{} {
	output := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &output)
	if err != nil {
		klog.Error(err)
		return ""
	}
	return output
}

func DeserializeForType(b []byte, objType reflect.Type) (*interface{}, error) {
	var objMap map[string]*json.RawMessage
	var err error
	err = json.Unmarshal(b, &objMap)
	if err != nil {
		return nil, err
	}

	objPtr := reflect.New(objType)

	err = json.Unmarshal(*objMap["K8sResource"], &objPtr)

	if err != nil {
		return nil, err
	}
	return objPtr.Elem().Interface(), nil
}
