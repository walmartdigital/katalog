package utils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/walmartdigital/katalog/src/domain"
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

// DeserializeResourceArray ...
func DeserializeResourceArray(b []byte, objType reflect.Type) ([]*domain.Resource, error) {
	var objMapArray []map[string]*json.RawMessage

	err := json.Unmarshal(b, &objMapArray)

	if err != nil {
		return nil, err
	}

	output := make([]*domain.Resource, len(objMapArray))

	for i, m := range objMapArray {
		r, err := DeserializeResource(m, objType)
		if err != nil {
			fmt.Println(err)
		}
		output[i] = r
	}
	return output, nil
}

// DeserializeResource ...
func DeserializeResource(objMap map[string]*json.RawMessage, objType reflect.Type) (*domain.Resource, error) {
	obj := reflect.New(objType).Interface()
	err1 := json.Unmarshal(*objMap["K8sResource"], &obj)

	if err1 != nil {
		return nil, err1
	}

	d := new(domain.Resource)
	d.K8sResource = obj.(domain.K8sResource)
	return d, nil
}
