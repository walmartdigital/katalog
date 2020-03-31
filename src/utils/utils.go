package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
)

var log = logrus.New()

func init() {
	err := LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
}

// Serialize ...
func Serialize(input interface{}) string {
	serial, err := json.Marshal(input)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(serial)
}

// Deserialize ...
func Deserialize(input string) interface{} {
	output := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &output)
	if err != nil {
		log.Error(err)
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

// ContainersToString ...
func ContainersToString(containers map[string]string) string {
	result := ""
	for k, v := range containers {
		if result == "" {
			result = fmt.Sprintf("%s:%s", k, v)
		} else {
			result = result + "," + fmt.Sprintf("%s:%s", k, v)
		}
	}
	return result
}

// LogInit ...
func LogInit(log *logrus.Logger) error {
	log.Formatter = &logrus.JSONFormatter{}
	log.SetOutput(os.Stdout)

	if logLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch logLevel {
		case "DEBUG":
			log.Info("Setting log level to DEBUG")
			log.SetLevel(logrus.DebugLevel)

		case "WARN":
			log.Info("Setting log level to WARN")
			log.SetLevel(logrus.WarnLevel)

		case "INFO":
			log.Info("Setting log level to WARN")
			log.SetLevel(logrus.InfoLevel)

		case "ERROR":
			log.Info("Setting log level to ERROR")
			log.SetLevel(logrus.ErrorLevel)

		default:
			log.Info("Setting log level to ERROR")
			log.SetLevel(logrus.ErrorLevel)
		}
	}

	if logLocation, ok := os.LookupEnv("LOG_FILE"); ok {
		file, err := os.OpenFile(logLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Info("Failed to log to file, using default stderr")
			return err
		}
		log.Out = file
	}

	return nil
}
