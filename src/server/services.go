package server

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server/repositories"
	"github.com/walmartdigital/katalog/src/utils"
)

var log = logrus.New()

func init() {
	err := utils.LogInit(log)
	if err != nil {
		log.Fatal(err)
	}
}

// Service ...
type Service struct {
	resourcesRepository repositories.Repository
	metrics             Metrics
}

// MakeService ...
func MakeService(resourcesRepository repositories.Repository, metricsfactory MetricsFactory) Service {
	return Service{
		resourcesRepository: resourcesRepository,
		metrics:             metricsfactory.Create(),
	}
}

// CreateService ...
func (s *Service) CreateService(service domain.Service) error {
	log.WithFields(logrus.Fields{
		"id":   service.GetID(),
		"name": service.GetName(),
	}).Debug("Creating Service")

	resource := domain.Resource{
		K8sResource: &service,
	}

	errCreatingResource := s.resourcesRepository.CreateResource(resource)
	if errCreatingResource != nil {
		log.WithFields(logrus.Fields{
			"msg": errCreatingResource.Error(),
		}).Error("Creating Service")

		return errCreatingResource
	}

	return nil
}

// UpdateService ...
func (s *Service) UpdateService(service domain.Service) error {
	log.WithFields(logrus.Fields{
		"id":   service.GetID(),
		"name": service.GetName(),
	}).Debug("Updating Service")

	resource := domain.Resource{K8sResource: &service}

	_, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update service (id: %s)", resource.GetID())
		return err
	}

	return nil
}

// DeleteService ...
func (s *Service) DeleteService(id string) error {
	log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Deleting Service")

	_, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Errorf("You provided a non-existing ID: %s", id)
		return err
	}

	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		log.Errorf("Deleted service ID: %s", id)
		return err
	}

	return nil
}

// CreateDeployment ...
func (s *Service) CreateDeployment(deployment domain.Deployment) error {
	log.WithFields(logrus.Fields{
		"id":   deployment.GetID(),
		"name": deployment.GetName(),
	}).Debug("Creating Deployment")

	resource := domain.Resource{K8sResource: &deployment}

	err := s.resourcesRepository.CreateResource(resource)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Creating Deployment")
		return err
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":                    resource.GetID(),
		"k8s-resource-type":                  "Deployment",
		"k8s-resource-ns":                    resource.GetNamespace(),
		"k8s-resource-name":                  resource.GetName(),
		"k8s-resource-labels":                resource.GetLabels(),
		"k8s-resource-annotations":           resource.GetAnnotations(),
		"k8s-resource-generation":            resource.GetGeneration(),
		"k8s-pod-template.containers.images": utils.ContainersToString(deployment.GetContainers()),
		"k8s-action":                         "create",
	}).Infof("Deployment %s/%s created", resource.GetNamespace(), resource.GetName())
	s.metrics.IncrementCounter("createDeployment", resource.GetID(), resource.GetNamespace(), resource.GetName())

	return nil
}

// UpdateDeployment ...
func (s *Service) UpdateDeployment(deployment domain.Deployment) error {
	log.WithFields(logrus.Fields{
		"id":   deployment.GetID(),
		"name": deployment.GetName(),
	}).Debug("Updating Deployment")

	resource := domain.Resource{K8sResource: &deployment}

	result, err := s.resourcesRepository.UpdateResource(resource)
	if err != nil {
		log.Errorf("Error occurred trying to update deployment (id: %s)", resource.GetID())
		return err
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":                    resource.GetID(),
		"k8s-resource-type":                  "Deployment",
		"k8s-resource-ns":                    resource.GetNamespace(),
		"k8s-resource-name":                  resource.GetName(),
		"k8s-resource-labels":                resource.GetLabels(),
		"k8s-resource-annotations":           resource.GetAnnotations(),
		"k8s-resource-generation":            resource.GetGeneration(),
		"k8s-pod-template.containers.images": utils.ContainersToString(deployment.GetContainers()),
		"k8s-action":                         "update",
	}).Infof("Deployment %s/%s updated", resource.GetNamespace(), resource.GetName())

	if result != nil {
		s.metrics.IncrementCounter("updateDeployment", resource.GetID(), resource.GetNamespace(), resource.GetName())
	}

	return nil
}

// DeleteDeployment ...
func (s *Service) DeleteDeployment(id string) error {
	log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Deleting Deployment")

	res, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Errorf("You provided a non-existing ID: %s", id)
		return err
	}

	if res == nil {
		log.WithFields(logrus.Fields{
			"id": id,
		}).Error("Delete Deployment Resource is null")

		return errors.New("Delete Deployment Resource null:" + id)
	}

	rep := res.(domain.Resource)
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		log.Errorf("Deleted deployment ID: %s", id)
		return err
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":          rep.GetID(),
		"k8s-resource-type":        "Deployment",
		"k8s-resource-ns":          rep.GetNamespace(),
		"k8s-resource-name":        rep.GetName(),
		"k8s-resource-labels":      rep.GetLabels(),
		"k8s-resource-annotations": rep.GetAnnotations(),
		"k8s-resource-generation":  rep.GetGeneration(),
		"k8s-action":               "delete",
	}).Infof("Deployment %s/%s deleted", rep.GetNamespace(), rep.GetName())

	s.metrics.IncrementCounter("deleteDeployment", id, rep.GetNamespace(), rep.GetName())

	return nil
}

// CreateStatefulSet ...
func (s *Service) CreateStatefulSet(statefulset domain.StatefulSet) error {
	log.WithFields(logrus.Fields{
		"id":   statefulset.GetID(),
		"name": statefulset.GetName(),
	}).Debug("Creating Statefulset")

	resource := domain.Resource{K8sResource: &statefulset}

	err := s.resourcesRepository.CreateResource(resource)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Create StatefulSet")
		return err
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":                    resource.GetID(),
		"k8s-resource-type":                  "StatefulSet",
		"k8s-resource-ns":                    resource.GetNamespace(),
		"k8s-resource-name":                  resource.GetName(),
		"k8s-resource-labels":                resource.GetLabels(),
		"k8s-resource-annotations":           resource.GetAnnotations(),
		"k8s-resource-generation":            resource.GetGeneration(),
		"k8s-pod-template.containers.images": utils.ContainersToString(statefulset.GetContainers()),
		"k8s-action":                         "create",
	}).Infof("Statefulset %s/%s created", resource.GetNamespace(), resource.GetName())

	s.metrics.IncrementCounter("createStatefulSet", resource.GetID(), resource.GetNamespace(), resource.GetName())

	return nil
}

// UpdateStatefulSet ...
func (s *Service) UpdateStatefulSet(statefulset domain.StatefulSet) error {
	log.WithFields(logrus.Fields{
		"id":   statefulset.GetID(),
		"name": statefulset.GetName(),
	}).Debug("Updating Statefulset")

	resource := domain.Resource{K8sResource: &statefulset}

	result, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update resource (id: %s)", resource.GetID())
		return err
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":                    resource.GetID(),
		"k8s-resource-type":                  "StatefulSet",
		"k8s-resource-ns":                    resource.GetNamespace(),
		"k8s-resource-name":                  resource.GetName(),
		"k8s-resource-labels":                resource.GetLabels(),
		"k8s-resource-annotations":           resource.GetAnnotations(),
		"k8s-resource-generation":            resource.GetGeneration(),
		"k8s-pod-template.containers.images": utils.ContainersToString(statefulset.GetContainers()),
		"k8s-action":                         "update",
	}).Infof("Statefulset %s/%s updated", resource.GetNamespace(), resource.GetName())

	if result != nil {
		s.metrics.IncrementCounter("updateStatefulSet", resource.GetID(), resource.GetNamespace(), resource.GetName())
	}

	return nil
}

// DeleteStatefulSet ...
func (s *Service) DeleteStatefulSet(id string) error {
	log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Deleting Statefulset")

	res, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Error("You have to provide an ID")
		return err
	}

	if res == nil {
		log.WithFields(logrus.Fields{
			"id": id,
		}).Error("Delete StatefulSet Resource is null")

		return errors.New("Delete StatefulSet Resource null:" + id)
	}

	rep := res.(domain.Resource)
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		log.Error("deleted statefulset id:" + id)
		return err
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":          rep.GetID(),
		"k8s-resource-type":        "StatefulSet",
		"k8s-resource-ns":          rep.GetNamespace(),
		"k8s-resource-name":        rep.GetName(),
		"k8s-resource-labels":      rep.GetLabels(),
		"k8s-resource-annotations": rep.GetAnnotations(),
		"k8s-resource-generation":  rep.GetGeneration(),
		"k8s-action":               "delete",
	}).Infof("Statefulset %s/%s deleted", rep.GetNamespace(), rep.GetName())

	s.metrics.IncrementCounter("deleteStatefulSet", id, rep.GetNamespace(), rep.GetName())

	return nil
}
