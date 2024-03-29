package publishers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/avast/retry-go"
	"github.com/walmartdigital/katalog/domain"
)

type httpClient interface {
}

// HTTPPublisher ...
type HTTPPublisher struct {
	url   string
	retry func(retry.RetryableFunc, ...retry.Option) error
}

// BuildHTTPPublisher ...
func BuildHTTPPublisher(url string, retry func(retry.RetryableFunc, ...retry.Option) error) Publisher {
	return &HTTPPublisher{url: url, retry: retry}
}

// Check ...
func (c *HTTPPublisher) Check() bool {
	return true
}

// Publish ...
func (c *HTTPPublisher) Publish(obj interface{}) error {
	operation := obj.(domain.Operation)
	switch operation.Kind {
	case (domain.OperationTypeAdd):
		return c.retry(func() error {
			return c.post(operation.Resource)
		})
	case (domain.OperationTypeUpdate):
		return c.retry(func() error {
			return c.put(operation.Resource)
		})
	case (domain.OperationTypeDelete):
		return c.delete(operation.Resource)
	default:
		return errors.New("operation unknown")
	}
}

func (c *HTTPPublisher) post(resource domain.Resource) error {
	reqBodyBytes := new(bytes.Buffer)

	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
		service := resource.GetK8sResource().(*domain.Service)
		err := json.NewEncoder(reqBodyBytes).Encode(*service)
		if err != nil {
			log.Error("Error serializing HTTP request body")
			return err
		}
		req, _ := http.NewRequest(http.MethodPost, c.url+"/services/"+service.ID, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("post service failed")
		}
		defer res.Body.Close()
		return nil

	case reflect.TypeOf(new(domain.Deployment)):
		deployment := resource.GetK8sResource().(*domain.Deployment)
		err := json.NewEncoder(reqBodyBytes).Encode(*deployment)
		if err != nil {
			log.Error("Error serializing HTTP request body")
			return err
		}
		req, _ := http.NewRequest(http.MethodPost, c.url+"/deployments/"+deployment.ID, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("post deployment failed")
		}
		defer res.Body.Close()
		return nil

	case reflect.TypeOf(new(domain.StatefulSet)):
		statefulset := resource.GetK8sResource().(*domain.StatefulSet)
		err := json.NewEncoder(reqBodyBytes).Encode(*statefulset)
		if err != nil {
			log.Error("Error serializing HTTP request body")
			return err
		}
		req, _ := http.NewRequest(http.MethodPost, c.url+"/statefulsets/"+statefulset.ID, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("post statefulset failed")
		}
		defer res.Body.Close()
		return nil

	default:
		log.Errorf("Type %s not found", v)
	}

	return nil
}

func (c *HTTPPublisher) put(resource domain.Resource) error {
	reqBodyBytes := new(bytes.Buffer)

	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
		service := resource.GetK8sResource().(*domain.Service)
		err := json.NewEncoder(reqBodyBytes).Encode(*service)
		if err != nil {
			log.Error("Error serializing HTTP request body")
			return err
		}
		req, _ := http.NewRequest(http.MethodPut, c.url+"/services/"+service.ID, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("put service failed")
		}
		defer res.Body.Close()
		return nil

	case reflect.TypeOf(new(domain.Deployment)):
		deployment := resource.GetK8sResource().(*domain.Deployment)
		err := json.NewEncoder(reqBodyBytes).Encode(*deployment)
		if err != nil {
			log.Error("Error serializing HTTP request body")
			return err
		}
		req, _ := http.NewRequest(http.MethodPut, c.url+"/deployments/"+deployment.ID, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("put deployment failed")
		}
		defer res.Body.Close()
		return nil

	case reflect.TypeOf(new(domain.StatefulSet)):
		statefulset := resource.GetK8sResource().(*domain.StatefulSet)
		err := json.NewEncoder(reqBodyBytes).Encode(*statefulset)
		if err != nil {
			log.Error("Error serializing HTTP request body")
			return err
		}
		req, _ := http.NewRequest(http.MethodPut, c.url+"/statefulsets/"+statefulset.ID, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("put statefulset failed")
		}
		defer res.Body.Close()
		return nil

	default:
		log.Errorf("Type %s not found", v)
	}

	return nil
}

func (c *HTTPPublisher) delete(resource domain.Resource) error {
	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
		service := resource.GetK8sResource().(*domain.Service)
		req, _ := http.NewRequest(http.MethodDelete, c.url+"/services/"+service.ID, nil)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("delete service failed")
		}
		defer res.Body.Close()
		return nil

	case reflect.TypeOf(new(domain.Deployment)):
		deployment := resource.GetK8sResource().(*domain.Deployment)
		req, _ := http.NewRequest(http.MethodDelete, c.url+"/deployments/"+deployment.ID, nil)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("delete deployment failed")
		}
		defer res.Body.Close()
		return nil

	case reflect.TypeOf(new(domain.StatefulSet)):
		statefulset := resource.GetK8sResource().(*domain.StatefulSet)
		req, _ := http.NewRequest(http.MethodDelete, c.url+"/statefulsets/"+statefulset.ID, nil)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Error(err)
			return errors.New("delete statefulset failed")
		}
		defer res.Body.Close()
		return nil

	default:
		log.Errorf("Type %s not found", v)
	}
	return nil
}
