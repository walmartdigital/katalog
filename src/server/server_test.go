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
	fail        bool
}

func (r *fakeRepository) CreateResource(obj interface{}) error {
	resource := obj.(domain.Resource)

	if resource.GetID() == "" {
		return errors.New("")
	}
	r.persistence[resource.GetID()] = resource

	return nil
}

func (r *fakeRepository) UpdateResource(resource interface{}) (*domain.Resource, error) {
	res := resource.(domain.Resource)
	savedResource, ok := r.persistence[res.GetID()]
	if !ok {
		return nil, errors.New("ID does not exist")
	}
	if r.fail {
		return nil, errors.New("error trying to update on database")
	}
	sr := savedResource.(domain.Resource)
	if &sr != nil {
		if sr.GetGeneration() < res.GetGeneration() {
			r.persistence[res.GetID()] = res
			return &res, nil
		}
	}
	return nil, nil
}

func (r *fakeRepository) DeleteResource(obj interface{}) error {
	id := obj.(string)
	if id == "" {
		return errors.New("need to provide an ID")
	}
	if r.fail {
		return errors.New("error trying to delete on database")
	}
	r.persistence[id] = nil
	return nil
}

func (r *fakeRepository) GetAllResources() ([]interface{}, error) {
	resources := arraylist.New()
	if r.fail {
		return nil, errors.New("error trying to delete on database")
	}
	for _, resource := range r.persistence {
		resources.Add(resource)
	}
	return resources.Values(), nil
}

func (r *fakeRepository) GetResource(id string) (interface{}, error) {
	resource, ok := r.persistence[id]
	if !ok {
		return nil, errors.New("ID does not exist")
	}
	return resource, nil
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
		persistence   map[string]interface{}
		repository    fakeRepository
		router        fakeRouter
		httpServer    fakeHTTPServer
		katalogServer server.Server
	)

	BeforeEach(func() {
		persistence = make(map[string]interface{})
		repository = fakeRepository{persistence: persistence, fail: false}
		router = fakeRouter{}
		httpServer = fakeHTTPServer{}
		katalogServer = server.CreateServer(&httpServer, &repository, &router)
		katalogServer.Run()
	})

	It("should create a service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(service)
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodPut, "", body)
		rec := httptest.NewRecorder()

		routes[path+"@POST"](rec, req)

		b, _ := ioutil.ReadAll(rec.Body)
		var srv domain.Service
		json.Unmarshal(b, &srv)
		resource := repository.persistence[id].(domain.Resource)
		output := resource.GetK8sResource().(*domain.Service)
		Expect(srv).To(Equal(*output))
	})

	It("should update a service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id, Port: 8888, Generation: 1}
		resource := domain.Resource{K8sResource: &service}
		repository.persistence[id] = resource
		newService := domain.Service{ID: id, Port: 9999, Generation: 2}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(newService)
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodPut, "/services/"+id, body)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@PUT"](rec, req)

		newResource := domain.Resource{K8sResource: &newService}
		Expect(repository.persistence[id]).To(Equal(newResource))
	})

	It("should not update an non-existing service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id, Port: 8888, Generation: 1}
		resource := domain.Resource{K8sResource: &service}
		repository.persistence[id] = resource
		newService := domain.Service{ID: "22d080de-4138-446f-acd4-d4c13fe779ff", Port: 9999, Generation: 2}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(newService)
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodPut, "/services/"+id, body)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@PUT"](rec, req)

		newResource := domain.Resource{K8sResource: &newService}
		Expect(repository.persistence[id]).NotTo(Equal(newResource))
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

	It("should not delete a non-existent service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		repository.persistence[id] = service
		nonExistentID := "a9b313fc-4cf6-42b9-8fa0-f0e258ea5c06"
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/services/"+nonExistentID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": nonExistentID})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[nonExistentID]).To(BeNil())
	})

	It("should not delete service when repository error", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		resource := domain.Resource{K8sResource: &service}
		repository.persistence[id] = resource
		repository.fail = true
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/services/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[id]).NotTo(BeNil())
	})

	It("should list all services", func() {
		inputResource1 := domain.Resource{
			K8sResource: &domain.Service{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}
		inputResource2 := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}
		inputResource3 := domain.Resource{
			K8sResource: &domain.Service{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
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
			i := r.GetID()
			resource := repository.persistence[i].(domain.Resource)
			output := resource.GetK8sResource().(*domain.Service)
			Expect(r.GetK8sResource().(*domain.Service)).To(Equal(output))
		}
	})

	It("should not list any services", func() {
		inputResource := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}
		repository.persistence["22d080de-4138-446f-acd4-d4c13fe77912"] = inputResource
		path := "/services"
		rec := httptest.NewRecorder()
		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)
		//TODO: expect stuff
	})

	It("should count amount of services", func() {
		inputResource1 := domain.Resource{
			K8sResource: &domain.Service{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}
		inputResource2 := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}
		inputResource3 := domain.Resource{
			K8sResource: &domain.Service{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
		}
		repository.persistence["22d080de-4138-446f-acd4-d4c13fe77912"] = inputResource1
		repository.persistence["22d080de-ffff-446f-acd4-d4c13fe77912"] = inputResource2
		repository.persistence["22d080de-xxxx-446f-acd4-d4c13fe77912"] = inputResource3
		path := "/services/_count"
		rec := httptest.NewRecorder()

		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)
		var m struct{ Count int }
		json.Unmarshal(b, &m)
		Expect(m.Count).To(Equal(2))
	})

	It("should create a service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(deployment)
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodPost, "", body)
		rec := httptest.NewRecorder()

		routes[path+"@POST"](rec, req)

		b, _ := ioutil.ReadAll(rec.Body)
		var srv domain.Deployment
		json.Unmarshal(b, &srv)
		resource := repository.persistence[id].(domain.Resource)
		output := resource.GetK8sResource().(*domain.Deployment)
		Expect(srv).To(Equal(*output))
	})

	It("should update a deployment", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id, Generation: 1}
		resource := domain.Resource{K8sResource: &deployment}
		repository.persistence[id] = resource
		newDeployment := domain.Deployment{ID: id, Generation: 2}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(newDeployment)
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodPut, "/deployments/"+id, body)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@PUT"](rec, req)

		newResource := domain.Resource{K8sResource: &newDeployment}
		Expect(repository.persistence[id]).To(Equal(newResource))
	})

	It("should not update an non-existing deployment", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id, Generation: 1}
		resource := domain.Resource{K8sResource: &deployment}
		repository.persistence[id] = resource
		newDeployment := domain.Deployment{ID: "22d080de-4138-446f-acd4-d4c13fe779ff", Generation: 2}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(newDeployment)
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodPut, "/deployments/"+id, body)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@PUT"](rec, req)

		newResource := domain.Resource{K8sResource: &newDeployment}
		Expect(repository.persistence[id]).NotTo(Equal(newResource))
	})

	It("should delete a deployment", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id}
		resource := domain.Resource{K8sResource: &deployment}
		repository.persistence[id] = resource
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/deployments/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[id]).To(BeNil())
	})

	It("should not delete deployment an non-existing deployment", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id}
		resource := domain.Resource{K8sResource: &deployment}
		repository.persistence[id] = resource
		nonExistingID := "3c665874-bc02-4691-981a-73fed8a12562"
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/deployments/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": nonExistingID})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[nonExistingID]).To(BeNil())
	})

	It("should not delete deployment when repository error", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		deployment := domain.Deployment{ID: id}
		resource := domain.Resource{K8sResource: &deployment}
		repository.persistence[id] = resource
		repository.fail = true
		path := "/deployments/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/deployments/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path+"@DELETE"](rec, req)

		Expect(repository.persistence[id]).NotTo(BeNil())
	})

	It("should list all deployments", func() {
		inputResource1 := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}
		inputResource2 := domain.Resource{
			K8sResource: &domain.Service{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}
		inputResource3 := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
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
			i := r.GetID()
			resource := repository.persistence[i].(domain.Resource)
			output := resource.GetK8sResource().(*domain.Deployment)
			Expect(r.GetK8sResource().(*domain.Deployment)).To(Equal(output))
		}
	})

	It("should count amount of deployments", func() {
		inputResource1 := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-4138-446f-acd4-d4c13fe77912"},
		}
		inputResource2 := domain.Resource{
			K8sResource: &domain.Service{ID: "22d080de-ffff-446f-acd4-d4c13fe77912"},
		}
		inputResource3 := domain.Resource{
			K8sResource: &domain.Deployment{ID: "22d080de-xxxx-446f-acd4-d4c13fe77912"},
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

	AfterEach(func() {
		katalogServer.DestroyMetrics()
	})
})
