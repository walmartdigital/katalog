package publishers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	retry "github.com/avast/retry-go"
	"github.com/golang/glog"
	"github.com/walmartdigital/katalog/src/domain"
)

type httpClient interface {
}

// HTTPPublisher ...
type HTTPPublisher struct {
	url string
}

// BuildHTTPPublisher ...
func BuildHTTPPublisher(url string) Publisher {
	return &HTTPPublisher{url: url}
}

// Publish ...
func (c *HTTPPublisher) Publish(obj interface{}) error {
	operation := obj.(domain.Operation)
	switch operation.Kind {
	case (domain.OperationTypeAdd):
		return c.put(operation.Service)
	case (domain.OperationTypeUpdate):
		return retry.Do(func() error {
			c.put(operation.Service)
			return errors.New("some error")
		})
	case (domain.OperationTypeDelete):
		return c.delete(operation.Service)
	default:
		return errors.New("operation unknown")
	}
}

func (c *HTTPPublisher) put(service domain.Service) error {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(service)
	req, _ := http.NewRequest(http.MethodPut, c.url+"/services/"+service.ID, reqBodyBytes)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		glog.Error(err)
		return errors.New("put service failed")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	glog.Info(string(body))
	return nil
}

func (c *HTTPPublisher) delete(service domain.Service) error {
	req, _ := http.NewRequest(http.MethodDelete, c.url+"/services/"+service.ID, nil)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		glog.Error(err)
		return errors.New("delete service failed")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	glog.Info(string(body))
	return nil
}
