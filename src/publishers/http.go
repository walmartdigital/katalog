package publishers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/seadiaz/katalog/src/domain"
)

// HTTPPublisher ...
type HTTPPublisher struct {
	url string
}

// CreateHTTPPublisher ...
func CreateHTTPPublisher(url string) Publisher {
	return &HTTPPublisher{url: url}
}

// Publish ...
func (c *HTTPPublisher) Publish(obj interface{}) {
	operation := obj.(domain.Operation)
	switch operation.Kind {
	case (domain.OperationTypeAdd):
		c.put(operation.Service)
	case (domain.OperationTypeDelete):
	}
}

func (c *HTTPPublisher) put(service domain.Service) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(service)
	glog.Info(reqBodyBytes)
	req, _ := http.NewRequest(http.MethodPut, c.url+"/services/"+service.ID, reqBodyBytes)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Error(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	glog.Info(string(body))

}
