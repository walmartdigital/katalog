package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gorilla/mux"
	"github.com/walmartdigital/katalog/src/domain"
)

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)
	resource := domain.Resource{
		Type:   "Service",
		Object: service,
	}
	s.resourcesRepository.CreateResource(resource)
	json.NewEncoder(w).Encode(service)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s.resourcesRepository.DeleteResource(id)
	fmt.Fprintf(w, "deleted service id: %s", id)
}

func getResourcesByType(resourceType string, resources []interface{}) []interface{} {
	list := arraylist.New()
	for _, r := range resources {
		res := r.(domain.Resource)
		if string(res.Type) == resourceType {
			list.Add(r)
		}
	}
	return list.Values()
}

func (s *Server) getAllServices(w http.ResponseWriter, r *http.Request) {
	resources := s.resourcesRepository.GetAllResources()
	services := getResourcesByType("Service", resources)
	json.NewEncoder(w).Encode(services)
}

func (s *Server) countServices(w http.ResponseWriter, r *http.Request) {
	services := s.resourcesRepository.GetAllResources()
	json.NewEncoder(w).Encode(struct{ Count int }{len(services)})
}

func (s *Server) createDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment domain.Deployment
	json.NewDecoder(r.Body).Decode(&deployment)
	resource := domain.Resource{
		Type:   "Deployment",
		Object: deployment,
	}
	s.resourcesRepository.CreateResource(resource)
	json.NewEncoder(w).Encode(deployment)
}

func (s *Server) deleteDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s.resourcesRepository.DeleteResource(id)
	fmt.Fprintf(w, "deleted deployment id: %s", id)
}

func (s *Server) getAllDeployments(w http.ResponseWriter, r *http.Request) {
	resources := s.resourcesRepository.GetAllResources()
	deployments := getResourcesByType("Deployment", resources)
	json.NewEncoder(w).Encode(deployments)
}

func (s *Server) countDeployments(w http.ResponseWriter, r *http.Request) {
	deployments := s.resourcesRepository.GetAllResources()
	json.NewEncoder(w).Encode(struct{ Count int }{len(deployments)})
}
