package kafka

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
)

// CreateService ...
func (c *Consumer) CreateService(body string) {
	var service domain.Service
	errDecoding := json.Unmarshal([]byte(body), &service)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.CreateService(service)
	if err != nil {
		logrus.Fatal(err)
	}
}

// UpdateService ...
func (c *Consumer) UpdateService(body string) {
	var service domain.Service
	errDecoding := json.Unmarshal([]byte(body), &service)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.UpdateService(service)
	if err != nil {
		logrus.Fatal(err)
	}
}

// DeleteService ...
func (c *Consumer) DeleteService(id string) {
	err := c.service.DeleteService(id)
	if err != nil {
		logrus.Fatal(err)
	}
}

// CreateDeployment ...
func (c *Consumer) CreateDeployment(body string) {
	var deployment domain.Deployment
	errDecoding := json.Unmarshal([]byte(body), &deployment)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.CreateDeployment(deployment)
	if err != nil {
		logrus.Fatal(err)
	}
}

// UpdateDeployment ...
func (c *Consumer) UpdateDeployment(body string) {
	var deployment domain.Deployment
	errDecoding := json.Unmarshal([]byte(body), &deployment)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.UpdateDeployment(deployment)
	if err != nil {
		logrus.Fatal(err)
	}
}

// DeleteDeployment ...
func (c *Consumer) DeleteDeployment(id string) {
	err := c.service.DeleteDeployment(id)
	if err != nil {
		logrus.Fatal(err)
	}
}

// CreateStatefulSet ...
func (c *Consumer) CreateStatefulSet(body string) {
	var statefulset domain.StatefulSet
	errDecoding := json.Unmarshal([]byte(body), &statefulset)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.CreateStatefulSet(statefulset)
	if err != nil {
		logrus.Fatal(err)
	}
}

// UpdateStatefulSet ...
func (c *Consumer) UpdateStatefulSet(body string) {
	var statefulset domain.StatefulSet
	errDecoding := json.Unmarshal([]byte(body), &statefulset)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.UpdateStatefulSet(statefulset)
	if err != nil {
		logrus.Fatal(err)
	}
}

// DeleteStatefulSet ...
func (c *Consumer) DeleteStatefulSet(id string) {
	err := c.service.DeleteStatefulSet(id)
	if err != nil {
		logrus.Fatal(err)
	}
}
