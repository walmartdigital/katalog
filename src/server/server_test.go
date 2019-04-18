package server_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/server"
	"github.com/walmartdigital/katalog/src/utils"
)

type dummyRepository struct{}

func (r *dummyRepository) CreateService(obj interface{}) {}
func (r *dummyRepository) DeleteService(obj interface{}) {}
func (r *dummyRepository) GetAllServices() []interface{} {
	return nil
}

type dummyRouter struct{}

func (r *dummyRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("***")
}

var router = mux.NewRouter().StrictSlash(true)
var routes = make(map[string]func(http.ResponseWriter, *http.Request))

func (r *dummyRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	fmt.Println(path)
	routes[path] = f
	return nil
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "server")
}

var _ = Describe("run server", func() {
	It("should receive all services", func() {
		path := "/services"
		req, _ := http.NewRequest("GET", "localhost:10000"+path, nil)
		recorder := httptest.NewRecorder()
		repository := &dummyRepository{}
		router := &dummyRouter{}
		server := server.CreateServer(repository, router)
		routes[path](recorder, req)
		res := recorder.Result()
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
		service := string(bytes.TrimSpace(b))

		server.Run()

		Expect(utils.Deserialize(service)).To(Equal())

	})
})
