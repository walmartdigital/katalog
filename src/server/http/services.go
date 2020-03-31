package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gorilla/mux"
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
		log.Fatal(errDecoding)
	}

	s.service.CreateService(service)

	errEncoding := json.NewEncoder(w).Encode(service)
	if errEncoding != nil {
		log.Fatal(errEncoding)
	}
}

// UpdateService ...
func (s *Server) UpdateService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	errDecoding := json.NewDecoder(r.Body).Decode(&service)
	if errDecoding != nil {
		log.Fatal(errDecoding)
	}

	s.service.UpdateService(service)

	errEncoding := json.NewEncoder(w).Encode(service)
	if errEncoding != nil {
		log.Fatal(errEncoding)
	}
}

// DeleteService ...
func (s *Server) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	s.service.DeleteService(id)

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
		log.Fatal(errEncoding)
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
		log.Fatal(errEncode)
	}
}

// CreateDeployment ...
func (s *Server) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	json.NewDecoder(r.Body).Decode(&deployment)

	s.service.CreateDeployment(deployment)
}

// UpdateDeployment ...
func (s *Server) UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	json.NewDecoder(r.Body).Decode(&deployment)

	s.service.UpdateDeployment(deployment)
}

// DeleteDeployment ...
func (s *Server) DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	s.service.DeleteDeployment(id)

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

// CreateStatefulSet ...
func (s *Server) CreateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulset domain.StatefulSet
	json.NewDecoder(r.Body).Decode(&statefulset)

	s.service.CreateStatefulSet(statefulset)
}

// UpdateStatefulSet ...
func (s *Server) UpdateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulset domain.StatefulSet
	json.NewDecoder(r.Body).Decode(&statefulset)

	s.service.UpdateStatefulSet(statefulset)
}

// DeleteStatefulSet ...
func (s *Server) DeleteStatefulSet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	s.service.DeleteStatefulSet(id)

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
