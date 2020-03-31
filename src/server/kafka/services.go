package kafka

import (
	"encoding/json"

	"github.com/walmartdigital/katalog/src/domain"
)

// CreateService ...
func (c *Consumer) CreateService(body string) {
	var service domain.Service
	errDecoding := json.Unmarshal([]byte(body), &service)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	c.service.CreateService(service)
}

// UpdateService ...
func (c *Consumer) UpdateService(body string) {
	var service domain.Service
	errDecoding := json.Unmarshal([]byte(body), &service)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	c.service.UpdateService(service)
}

// DeleteService ...
func (c *Consumer) DeleteService(id string) {
	c.service.DeleteService(id)
}

// CreateDeployment ...
func (c *Consumer) CreateDeployment(body string) {
	var deployment domain.Deployment
	errDecoding := json.Unmarshal([]byte(body), &deployment)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	c.service.CreateDeployment(deployment)
}

// UpdateDeployment ...
func (c *Consumer) UpdateDeployment(body string) {
	var deployment domain.Deployment
	errDecoding := json.Unmarshal([]byte(body), &deployment)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	c.service.UpdateDeployment(deployment)
}

// DeleteDeployment ...
func (c *Consumer) DeleteDeployment(id string) {
	c.service.DeleteDeployment(id)
}

// CreateStatefulSet ...
func (c *Consumer) CreateStatefulSet(body string) {
	var statefulset domain.StatefulSet
	errDecoding := json.Unmarshal([]byte(body), &statefulset)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	c.service.CreateStatefulSet(statefulset)
}

// UpdateStatefulSet ...
func (c *Consumer) UpdateStatefulSet(body string) {
	var statefulset domain.StatefulSet
	errDecoding := json.Unmarshal([]byte(body), &statefulset)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	c.service.UpdateStatefulSet(statefulset)
}

// DeleteStatefulSet ...
func (c *Consumer) DeleteStatefulSet(id string) {
	c.service.DeleteStatefulSet(id)
}
