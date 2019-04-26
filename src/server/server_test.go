package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gorilla/mux"
	"github.com/maxcnunes/httpfake"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/server"
)

const kind = "services"

type fakeRepository struct {
	persistence map[string]interface{}
}

func (r *fakeRepository) CreateService(obj interface{}) {
	r.persistence[kind] = obj
}

func (r *fakeRepository) DeleteService(obj interface{}) {
	r.persistence[kind] = nil
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
		path        string
		persistence map[string]interface{}
		repository  fakeRepository
		router      fakeRouter
		httpServer  fakeHTTPServer
	)

	BeforeEach(func() {
		path = "/services"
		persistence = make(map[string]interface{})
		repository = fakeRepository{persistence: persistence}
		router = fakeRouter{}
		httpServer = fakeHTTPServer{}
		server := server.CreateServer(&httpServer, &repository, &router)
		server.Run()
	})

	It("should list all services", func() {
		rec := httptest.NewRecorder()
		repository.persistence[kind] = struct{ ID string }{"c3642c80-b6cd-4e5c-ab9a-64ce340a8c83"}

		routes[path](rec, nil)

		res := rec.Result()
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())
		var services []struct{ ID string }
		Expect(json.Unmarshal(body, &services)).To(BeNil())
		Expect(services[0]).To(Equal(repository.persistence[kind]))
	})

	It("should create a service", func() {
		fakeService := httpfake.New()
		defer fakeService.Server.Close()
		serviceID := "22d080de-4138-446f-acd4-d4c13fe77912"
		fakeService.NewHandler().
			Put("/services/" + serviceID).
			Reply(200).
			BodyString(`{"ID": "` + serviceID + `"}`)
		url := fakeService.ResolveURL("")
		path = fmt.Sprintf("%s/%s", path, serviceID)
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(struct{ ID string }{ID: serviceID})
		req, _ := http.NewRequest(http.MethodPut, url, reqBodyBytes)
		req.Header.Add("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		routes[path](rec, req)

		Expect(repository.persistence[kind]).To(Equal(struct{ ID string }{ID: serviceID}))
	})
})
