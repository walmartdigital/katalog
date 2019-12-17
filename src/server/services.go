package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
)

func (s *Server) getResourcesByType(resource domain.Resource) ([]interface{}, error) {
	resources, err := s.resourcesRepository.GetAllResources()
	if err != nil {
		return nil, err
	}

	list := arraylist.New()
	for _, r := range resources {
		res := r.(domain.Resource)
		if res.GetType() == resource.GetType() {
			list.Add(r)
		}
	}
	return list.Values(), nil
}

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	resource := domain.Resource{
		K8sResource: &service,
	}
	s.resourcesRepository.CreateResource(resource)
	json.NewEncoder(w).Encode(service)
}

func (s *Server) updateService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	resource := domain.Resource{K8sResource: &service}
	_, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update service (id: %s)", resource.GetID())
		return
	}

	json.NewEncoder(w).Encode(service)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		fmt.Fprintf(w, "You provided a non-existing ID: %s", id)
		log.Errorf("You provided a non-existing ID: %s", id)
		return
	}
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		fmt.Fprintf(w, "Deleted service ID: %s", id)
		log.Errorf("Deleted service ID: %s", id)
		return
	}
	fmt.Fprintf(w, "deleted service id: %s", id)
}

func (s *Server) getAllServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Service{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) countServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Service{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	json.NewEncoder(w).Encode(struct{ Count int }{len(services)})
}

func (s *Server) createDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	json.NewDecoder(r.Body).Decode(&deployment)
	resource := domain.Resource{K8sResource: &deployment}
	s.resourcesRepository.CreateResource(resource)

	log.WithFields(logrus.Fields{
		"k8s-resource-id":           resource.GetID(),
		"k8s-resource-type":         "Deployment",
		"k8s-resource-ns":           resource.GetNamespace(),
		"k8s-resource-name":         resource.GetName(),
		"k8s-resource-labels":       resource.GetLabels(),
		"k8s-deployment-containers": deployment.GetContainers(),
		"k8s-action":                "create",
	}).Infof("Deployment %s/%s created", resource.GetNamespace(), resource.GetName())

	(*s.metrics)["createDeployment"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
	json.NewEncoder(w).Encode(deployment)
}

func (s *Server) updateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	json.NewDecoder(r.Body).Decode(&deployment)
	resource := domain.Resource{K8sResource: &deployment}
	result, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update deployment (id: %s)", resource.GetID())
		return
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":           resource.GetID(),
		"k8s-resource-type":         "Deployment",
		"k8s-resource-ns":           resource.GetNamespace(),
		"k8s-resource-name":         resource.GetName(),
		"k8s-resource-labels":       resource.GetLabels(),
		"k8s-deployment-containers": deployment.GetContainers(),
		"k8s-action":                "update",
	}).Infof("Deployment %s/%s updated", resource.GetNamespace(), resource.GetName())

	if result != nil {
		(*s.metrics)["updateDeployment"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
	}
	json.NewEncoder(w).Encode(deployment)
}

func (s *Server) deleteDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	res, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		fmt.Fprintf(w, "You provided a non-existing ID: %s", id)
		log.Errorf("You provided a non-existing ID: %s", id)
		return
	}
	rep := res.(domain.Resource)
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		fmt.Fprintf(w, "Deleted deployment ID: %s", id)
		log.Errorf("Deleted deployment ID: %s", id)
		return
	}
	fmt.Fprintf(w, "deleted deployment id: %s", id)

	log.WithFields(logrus.Fields{
		"k8s-resource-id":     rep.GetID(),
		"k8s-resource-type":   "Deployment",
		"k8s-resource-ns":     rep.GetNamespace(),
		"k8s-resource-name":   rep.GetName(),
		"k8s-resource-labels": rep.GetLabels(),
		"k8s-action":          "delete",
	}).Infof("Deployment %s/%s deleted", rep.GetNamespace(), rep.GetName())

	(*s.metrics)["deleteDeployment"].(*prometheus.CounterVec).WithLabelValues(id, rep.GetNamespace(), rep.GetName()).Inc()
}

func (s *Server) getAllDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Deployment{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

func (s *Server) countDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Deployment{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	json.NewEncoder(w).Encode(struct{ Count int }{len(deployments)})
}

func (s *Server) createStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulset domain.StatefulSet
	json.NewDecoder(r.Body).Decode(&statefulset)
	resource := domain.Resource{K8sResource: &statefulset}
	s.resourcesRepository.CreateResource(resource)

	log.WithFields(logrus.Fields{
		"k8s-resource-id":           resource.GetID(),
		"k8s-resource-type":         "StatefulSet",
		"k8s-resource-ns":           resource.GetNamespace(),
		"k8s-resource-name":         resource.GetName(),
		"k8s-resource-labels":       resource.GetLabels(),
		"k8s-deployment-containers": statefulset.GetContainers(),
		"k8s-action":                "create",
	}).Infof("Statefulset %s/%s created", resource.GetNamespace(), resource.GetName())

	(*s.metrics)["createStatefulSet"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
	json.NewEncoder(w).Encode(statefulset)
}

func (s *Server) updateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulset domain.StatefulSet
	json.NewDecoder(r.Body).Decode(&statefulset)
	resource := domain.Resource{K8sResource: &statefulset}
	result, err := s.resourcesRepository.UpdateResource(resource)

	if err != nil {
		log.Errorf("Error occurred trying to update resource (id: %s)", resource.GetID())
		return
	}

	log.WithFields(logrus.Fields{
		"k8s-resource-id":           resource.GetID(),
		"k8s-resource-type":         "StatefulSet",
		"k8s-resource-ns":           resource.GetNamespace(),
		"k8s-resource-name":         resource.GetName(),
		"k8s-resource-labels":       resource.GetLabels(),
		"k8s-deployment-containers": statefulset.GetContainers(),
		"k8s-action":                "update",
	}).Infof("Statefulset %s/%s updated", resource.GetNamespace(), resource.GetName())

	if result != nil {
		(*s.metrics)["updateStatefulSet"].(*prometheus.CounterVec).WithLabelValues(resource.GetID(), resource.GetNamespace(), resource.GetName()).Inc()
	}
	json.NewEncoder(w).Encode(statefulset)
}

func (s *Server) deleteStatefulSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	res, err := s.resourcesRepository.GetResource(id)
	if err != nil {
		log.Error("You have to provide an ID")
		return
	}
	rep := res.(domain.Resource)
	err = s.resourcesRepository.DeleteResource(id)
	if err != nil {
		fmt.Fprintf(w, "deleted statefulset id: %s", id)
		return
	}
	fmt.Fprintf(w, "deleted statefulset id: %s", id)

	log.WithFields(logrus.Fields{
		"k8s-resource-id":     rep.GetID(),
		"k8s-resource-type":   "StatefulSet",
		"k8s-resource-ns":     rep.GetNamespace(),
		"k8s-resource-name":   rep.GetName(),
		"k8s-resource-labels": rep.GetLabels(),
		"k8s-action":          "delete",
	}).Infof("Statefulset %s/%s deleted", rep.GetNamespace(), rep.GetName())

	(*s.metrics)["deleteStatefulSet"].(*prometheus.CounterVec).WithLabelValues(id, rep.GetNamespace(), rep.GetName()).Inc()
}

func (s *Server) getAllStatefulSets(w http.ResponseWriter, r *http.Request) {
	statefulsets, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.StatefulSet{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statefulsets)
}

func (s *Server) countStatefulSets(w http.ResponseWriter, r *http.Request) {
	statefulsets, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.StatefulSet{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	json.NewEncoder(w).Encode(struct{ Count int }{len(statefulsets)})
}
