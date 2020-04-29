package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gorilla/mux"
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

// CreateService ...
func (s *Server) CreateService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	errDecoding := json.NewDecoder(r.Body).Decode(&service)
	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Error("Deserializing Service")
	}

	errCreating := s.service.CreateService(service)
	if errCreating != nil {
		log.WithFields(logrus.Fields{
			"msg": errCreating.Error(),
		}).Error("Creating Service")
	}

	errEncoding := json.NewEncoder(w).Encode(service)
	if errEncoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncoding.Error(),
		}).Error("Encoding Service")
	}
}

// UpdateService ...
func (s *Server) UpdateService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	errDecoding := json.NewDecoder(r.Body).Decode(&service)

	if errDecoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errDecoding.Error(),
		}).Error("Deserializing Service")
	}

	errUpdating := s.service.UpdateService(service)
	if errUpdating != nil {
		log.WithFields(logrus.Fields{
			"msg": errUpdating.Error(),
		}).Error("Updating Service")
	}

	errEncoding := json.NewEncoder(w).Encode(service)
	if errEncoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncoding.Error(),
		}).Error("Encoding Service")
	}
}

// DeleteService ...
func (s *Server) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.service.DeleteService(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deleting Service")
	}

	fmt.Fprintf(w, "service id: %s", id)
}

func (s *Server) getAllServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Service{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	errEncoding := json.NewEncoder(w).Encode(services)
	if errEncoding != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncoding.Error(),
		}).Error("Getting all services")
	}
}

func (s *Server) countServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Service{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	errEncode := json.NewEncoder(w).Encode(struct{ Count int }{len(services)})
	if errEncode != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncode.Error(),
		}).Error("Counting all services")
	}
}

// CreateDeployment ...
func (s *Server) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	err := json.NewDecoder(r.Body).Decode(&deployment)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deserializing Deployment")
	}

	err = s.service.CreateDeployment(deployment)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Creating Deployment")
	}

	errEncode := json.NewEncoder(w).Encode(deployment)
	if errEncode != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncode.Error(),
		}).Error("Encoding Create Deployment")
	}
}

// UpdateDeployment ...
func (s *Server) UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	err := json.NewDecoder(r.Body).Decode(&deployment)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deserializing Deployment")
	}

	err = s.service.UpdateDeployment(deployment)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Updating Deployment")
	}

	errEncode := json.NewEncoder(w).Encode(deployment)
	if errEncode != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncode.Error(),
		}).Error("Encoding Update Deployment")
	}
}

// DeleteDeployment ...
func (s *Server) DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.service.DeleteDeployment(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deleting Deployment")
	}

	fmt.Fprintf(w, "deployment id: %s", id)
}

func (s *Server) getAllDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Deployment{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deployments)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Getting all Deployments")
	}
}

func (s *Server) countDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.Deployment{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	err = json.NewEncoder(w).Encode(struct{ Count int }{len(deployments)})
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Counting all Deployments")
	}
}

// CreateStatefulSet ...
func (s *Server) CreateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulset domain.StatefulSet
	err := json.NewDecoder(r.Body).Decode(&statefulset)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deserializing StatefulSet")
	}

	err = s.service.CreateStatefulSet(statefulset)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Creating StatefulSet")
	}

	errEncode := json.NewEncoder(w).Encode(statefulset)
	if errEncode != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncode.Error(),
		}).Error("Encoding Create StatefulSet")
	}
}

// UpdateStatefulSet ...
func (s *Server) UpdateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulset domain.StatefulSet
	err := json.NewDecoder(r.Body).Decode(&statefulset)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deserializing StatefulSet")
	}

	err = s.service.UpdateStatefulSet(statefulset)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Updating StatefulSet")
	}

	errEncode := json.NewEncoder(w).Encode(statefulset)
	if errEncode != nil {
		log.WithFields(logrus.Fields{
			"msg": errEncode.Error(),
		}).Error("Encoding Update StatefulSet")
	}
}

// DeleteStatefulSet ...
func (s *Server) DeleteStatefulSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.service.DeleteStatefulSet(id)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Deleting StatefulSet")
	}

	fmt.Fprintf(w, "deployment id: %s", id)
}

func (s *Server) getAllStatefulSets(w http.ResponseWriter, r *http.Request) {
	statefulsets, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.StatefulSet{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(statefulsets)
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Getting All StatefulSet")
	}
}

func (s *Server) countStatefulSets(w http.ResponseWriter, r *http.Request) {
	statefulsets, err := s.getResourcesByType(domain.Resource{K8sResource: &domain.StatefulSet{}})
	if err != nil {
		fmt.Fprint(w, "Resource not found")
		log.Error("Resource not found")
		return
	}
	err = json.NewEncoder(w).Encode(struct{ Count int }{len(statefulsets)})
	if err != nil {
		log.WithFields(logrus.Fields{
			"msg": err.Error(),
		}).Error("Counting StatefulSet")
	}
}
