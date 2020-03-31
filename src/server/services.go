package server

import (
	"github.com/prometheus/client_golang/prometheus"
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
	metrics             *map[string]interface{}
}

// MakeService ...
func MakeService(resourcesRepository repositories.Repository, metrics *map[string]interface{}) Service {
	return Service{
		resourcesRepository: resourcesRepository,
		metrics:             metrics,
	}
}

// CreateService ...
func (s *Service) CreateService(service domain.Service) {
	resource := domain.Resource{
		K8sResource: &service,
	}
	errCreatingResource := s.resourcesRepository.CreateResource(resource)
	if errCreatingResource != nil {
		log.Fatal(errCreatingResource)
	}
}

// UpdateService ...
func (s *Service) UpdateService(service domain.Service) {
	resource := domain.Resource{K8sResource: &service}
	_, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update service (id: %s)", resource.GetID())
		return
	}
}

// DeleteService ...
func (s *Service) DeleteService(id string) {
	_, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Errorf("You provided a non-existing ID: %s", id)
		return
	}
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		log.Errorf("Deleted service ID: %s", id)
		return
	}
}

// CreateDeployment ...
func (s *Service) CreateDeployment(deployment domain.Deployment) {
	resource := domain.Resource{K8sResource: &deployment}
	s.resourcesRepository.CreateResource(resource)

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

	(*s.metrics)["createDeployment"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
}

// UpdateDeployment ...
func (s *Service) UpdateDeployment(deployment domain.Deployment) {
	resource := domain.Resource{K8sResource: &deployment}
	result, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update deployment (id: %s)", resource.GetID())
		return
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
		(*s.metrics)["updateDeployment"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
	}
}

// DeleteDeployment ...
func (s *Service) DeleteDeployment(id string) {
	res, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Errorf("You provided a non-existing ID: %s", id)
		return
	}
	rep := res.(domain.Resource)
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		log.Errorf("Deleted deployment ID: %s", id)
		return
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

	(*s.metrics)["deleteDeployment"].(*prometheus.CounterVec).WithLabelValues(id, rep.GetNamespace(), rep.GetName()).Inc()
}

// CreateStatefulSet ...
func (s *Service) CreateStatefulSet(statefulset domain.StatefulSet) {
	resource := domain.Resource{K8sResource: &statefulset}
	s.resourcesRepository.CreateResource(resource)

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

	(*s.metrics)["createStatefulSet"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
}

// UpdateStatefulSet ...
func (s *Service) UpdateStatefulSet(statefulset domain.StatefulSet) {
	resource := domain.Resource{K8sResource: &statefulset}
	result, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update resource (id: %s)", resource.GetID())
		return
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
		(*s.metrics)["updateStatefulSet"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
	}
}

// DeleteStatefulSet ...
func (s *Service) DeleteStatefulSet(id string) {
	res, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Error("You have to provide an ID")
		return
	}
	rep := res.(domain.Resource)
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		log.Error("deleted statefulset id: %s", id)
		return
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

	(*s.metrics)["deleteStatefulSet"].(*prometheus.CounterVec).WithLabelValues(id, rep.GetNamespace(), rep.GetName()).Inc()
}
