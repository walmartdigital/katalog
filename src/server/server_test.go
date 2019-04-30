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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/server"
)

type fakeRepository struct {
	persistence map[string]interface{}
}

func (r *fakeRepository) CreateService(obj interface{}) error {
	service := obj.(domain.Service)
	if service.ID == "" {
		return errors.New("")
	}
	r.persistence[service.ID] = obj
	return nil
}

func (r *fakeRepository) DeleteService(obj interface{}) error {
	id := obj.(string)
	if id == "" {
		return errors.New("")
	}
	r.persistence[id] = nil
	return nil
}

func (r *fakeRepository) GetAllServices() []interface{} {
	services := arraylist.New()
	for _, service := range r.persistence {
		services.Add(service)
	}
	return services.Values()
}

type fakeRouter struct{}

func (r *fakeRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

type fakeHTTPServer struct{}

func (s *fakeHTTPServer) ListenAndServe() error {
	return nil
}

var routes = make(map[string]func(http.ResponseWriter, *http.Request))

func (r *fakeRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	routes[path] = f
	return &mux.Route{}
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
		path := "/services/{service}"
		req, _ := http.NewRequest(http.MethodPut, "", body)
		rec := httptest.NewRecorder()

		routes[path](rec, req)

		b, _ := ioutil.ReadAll(rec.Body)
		var srv domain.Service
		json.Unmarshal(b, &srv)
		Expect(srv).To(Equal(repository.persistence[id]))
	})

	It("should delete a service", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		repository.persistence[id] = service
		path := "/services/{id}"
		req, _ := http.NewRequest(http.MethodDelete, "/services/"+id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		rec := httptest.NewRecorder()

		routes[path](rec, req)

		Expect(repository.persistence[id]).To(BeNil())
	})

	It("should list all services", func() {
		id := "22d080de-4138-446f-acd4-d4c13fe77912"
		service := domain.Service{ID: id}
		repository.persistence[id] = service
		path := "/services"
		rec := httptest.NewRecorder()

		routes[path](rec, nil)

		b, _ := ioutil.ReadAll(rec.Body)
		json.Unmarshal(b, &service)
		Expect(service).To(Equal(repository.persistence[id]))
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
})
