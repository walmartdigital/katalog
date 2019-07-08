package server_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server"
)

type fakeRepository struct {
	persistence map[string]interface{}
}

func (r *fakeRepository) CreateResource(obj interface{}) error {
	resource := obj.(domain.Resource)
	if resource.Type == "Service" {
		if resource.Object.(domain.Service).ID == "" {
			return errors.New("")
		}
		r.persistence[resource.Object.(domain.Service).ID] = resource
	}
	if resource.Type == "Deployment" {
		if resource.Object.(domain.Deployment).ID == "" {
			return errors.New("")
		}
		r.persistence[resource.Object.(domain.Deployment).ID] = resource
	}
	return nil
}

func (r *fakeRepository) DeleteResource(obj interface{}) error {
	id := obj.(string)
	if id == "" {
		return errors.New("")
	}
	r.persistence[id] = nil
	return nil
}

func (r *fakeRepository) GetAllResources() []interface{} {
	resources := arraylist.New()
	for _, resource := range r.persistence {
		resources.Add(resource)
	}
	return resources.Values()
}

type fakeRouter struct{}

var routes = make(map[string]func(http.ResponseWriter, *http.Request))

func (r *fakeRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) server.Route {
	routes[path] = f
	return &fakeRoute{path}
}

type fakeRoute struct {
	path string
}

func (r *fakeRoute) Methods(methods ...string) server.Route {
	routes[r.path+"@"+methods[0]] = routes[r.path]
	return r
}

type fakeHTTPServer struct{}

func (s *fakeHTTPServer) ListenAndServe() error {
	return nil
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server")
}

var _ = Describe("run server", func() {
	var (
		persistence map[string]interface{}
		repository  fakeRepository
		router      fakeRouter
		httpServer  fakeHTTPServer
	)

	BeforeEach(func() {
		persistence = make(map[string]interface{})
		repository = fakeRepository{persistence: persistence}
		router = fakeRouter{}
		httpServer = fakeHTTPServer{}
		server := server.CreateServer(&httpServer, &repository, &router)
		server.Run()
	})

	It("should create a service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(service)
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodPut, "", body)
		rec := httptest.NewRecorder()

		routes[path+"@PUT"](rec, req)

		b, _ := ioutil.ReadAll(rec.Body)
		var srv domain.Service
		json.Unmarshal(b, &srv)
		resource := repository.persistence[id].(domain.Resource)
		output := resource.Object.(domain.Service)
		Expect(srv).To(Equal(output))
	})

	It("should delete a service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		repository.persistence[id] = service
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/services/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[id]).To(BeNil())
	})

	It("should list all services", func() {
		inputResource1 := domain.Resource{
			Type:   "Service",
			Object: domain.Service{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}
		inputResource2 := domain.Resource{
			Type:   "Service",
			Object: domain.Service{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}
		inputResource3 := domain.Resource{
			Type:   "Service",
			Object: domain.Service{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
		}

		repository.persistence["22d080de-4138-446f-acd4-d4c13fe77912"] = inputResource1
		repository.persistence["22d080de-ffff-446f-acd4-d4c13fe77912"] = inputResource2
		repository.persistence["22d080de-xxxx-446f-acd4-d4c13fe77912"] = inputResource3

		path := "/services"
		rec := httptest.NewRecorder()

		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)

		var resources []domain.Resource
		var d map[string]interface{}
		json.Unmarshal(b, &d)
		mapstructure.Decode(resources, &d)

		for _, r := range resources {
			i := r.Object.(domain.Service).ID
			resource := repository.persistence[i].(domain.Resource)
			output := resource.Object.(domain.Service)
			Expect(r.Object.(domain.Service)).To(Equal(output))
		}
	})

	It("should count amount of services", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		repository.persistence[id] = service
		path := "/services/_count"
		rec := httptest.NewRecorder()

		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)
		var m struct{ Count int }
		json.Unmarshal(b, &m)
		Expect(m.Count).To(Equal(1))
	})

	It("should create a deployment", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(deployment)
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodPut, "", body)
		rec := httptest.NewRecorder()

		routes[path+"@PUT"](rec, req)

		b, _ := ioutil.ReadAll(rec.Body)
		var srv domain.Deployment
		json.Unmarshal(b, &srv)
		resource := repository.persistence[id].(domain.Resource)
		output := resource.Object.(domain.Deployment)
		Expect(srv).To(Equal(output))
	})

	It("should delete a deployment", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id}
		repository.persistence[id] = deployment
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/deployments/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[id]).To(BeNil())
	})

	It("should list all deployments", func() {
		inputResource1 := domain.Resource{
			Type:   "Deployment",
			Object: domain.Deployment{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}

		inputResource2 := domain.Resource{
			Type:   "Service",
			Object: domain.Service{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}

		inputResource3 := domain.Resource{
			Type:   "Deployment",
			Object: domain.Deployment{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
		}

		repository.persistence["22d080de-4138-446f-acd4-d4c13fe77912"] = inputResource1
		repository.persistence["22d080de-ffff-446f-acd4-d4c13fe77912"] = inputResource2
		repository.persistence["22d080de-xxxx-446f-acd4-d4c13fe77912"] = inputResource3

		path := "/deployments"
		rec := httptest.NewRecorder()

		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)

		var resources []domain.Resource
		var d map[string]interface{}
		json.Unmarshal(b, &d)
		mapstructure.Decode(resources, &d)

		for _, r := range resources {
			i := r.Object.(domain.Deployment).ID
			resource := repository.persistence[i].(domain.Resource)
			output := resource.Object.(domain.Deployment)
			Expect(r.Object.(domain.Deployment)).To(Equal(output))
		}
	})

	It("should count amount of deployments", func() {
		inputResource1 := domain.Resource{
			Type:   "Deployment",
			Object: domain.Deployment{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}

		inputResource2 := domain.Resource{
			Type:   "Service",
			Object: domain.Service{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}

		inputResource3 := domain.Resource{
			Type:   "Deployment",
			Object: domain.Deployment{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
		}

		repository.persistence["22d080de-4138-446f-acd4-d4c13fe77912"] = inputResource1
		repository.persistence["22d080de-ffff-446f-acd4-d4c13fe77912"] = inputResource2
		repository.persistence["22d080de-xxxx-446f-acd4-d4c13fe77912"] = inputResource3

		path := "/deployments/_count"
		rec := httptest.NewRecorder()

		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)
		var m struct{ Count int }
		json.Unmarshal(b, &m)
		Expect(m.Count).To(Equal(2))
	})
})
