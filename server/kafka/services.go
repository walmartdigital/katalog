package kafka

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/domain"
)

// CreateService ...
func (c *Consumer) CreateService(body string) error {
	var service domain.Service
	errDecoding := json.Unmarshal([]byte(body), &service)
	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Debug("Deserializing Service")

		return errDecoding
	}

	err := c.service.CreateService(service)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Creating Service")

		return err
	}

	return nil
}

// UpdateService ...
func (c *Consumer) UpdateService(body string) error {
	var service domain.Service
	errDecoding := json.Unmarshal([]byte(body), &service)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	err := c.service.UpdateService(service)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Creating Service")

		return err
	}

	return nil
}

// DeleteService ...
func (c *Consumer) DeleteService(id string) error {
	err := c.service.DeleteService(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Deleting Service")

		return err
	}

	return nil
}

// CreateDeployment ...
func (c *Consumer) CreateDeployment(body string) error {
	var deployment domain.Deployment
	errDecoding := json.Unmarshal([]byte(body), &deployment)
	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Debug("Deserializing Deployment")

		return errDecoding
	}

	err := c.service.CreateDeployment(deployment)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Creating Deployment")

		return err
	}

	return nil
}

// UpdateDeployment ...
func (c *Consumer) UpdateDeployment(body string) error {
	var deployment domain.Deployment
	errDecoding := json.Unmarshal([]byte(body), &deployment)
	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Debug("Deserializing Deployment")

		return errDecoding
	}

	err := c.service.UpdateDeployment(deployment)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Updating Deployment")

		return err
	}

	return nil
}

// DeleteDeployment ...
func (c *Consumer) DeleteDeployment(id string) error {
	err := c.service.DeleteDeployment(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Deleting Deployment")

		return err
	}

	return nil
}

// CreateStatefulSet ...
func (c *Consumer) CreateStatefulSet(body string) error {
	var statefulset domain.StatefulSet
	errDecoding := json.Unmarshal([]byte(body), &statefulset)
	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Debug("Deserializing StatefulSet")

		return errDecoding
	}

	err := c.service.CreateStatefulSet(statefulset)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Creating StatefulSet")

		return err
	}

	return nil
}

// UpdateStatefulSet ...
func (c *Consumer) UpdateStatefulSet(body string) error {
	var statefulset domain.StatefulSet
	errDecoding := json.Unmarshal([]byte(body), &statefulset)
	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Debug("Deserializing StatefulSet")

		return errDecoding
	}

	err := c.service.UpdateStatefulSet(statefulset)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Updating StatefulSet")

		return err
	}

	return nil
}

// DeleteStatefulSet ...
func (c *Consumer) DeleteStatefulSet(id string) error {
	err := c.service.DeleteStatefulSet(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Debug("Creating StatefulSet")

		return err
	}

	return nil
}
