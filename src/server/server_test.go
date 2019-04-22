package server_test

import (
	"fmt"
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

type dummyRepository struct {
	services []domain.Service
}

func (r *dummyRepository) CreateService(obj interface{}) {}
func (r *dummyRepository) DeleteService(obj interface{}) {}
func (r *dummyRepository) GetAllServices() []interface{} {
	list := arraylist.New()
	return list.Values()
}

type dummyRouter struct{}

func (r *dummyRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("***")
}

// var router = mux.NewRouter().StrictSlash(true)
var routes = make(map[string]func(http.ResponseWriter, *http.Request))

func (r *dummyRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	fmt.Println(path)
	routes[path] = f
	route := &mux.Route{}
	return route
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server")
}

var _ = Describe("run server", func() {
	It("should receive status code 200", func() {
		repository := &dummyRepository{}
		router := &dummyRouter{}
		server := server.CreateServer(repository, router)
		server.Run()
		path := "/services"
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:10000"+path, nil)
		rec := httptest.NewRecorder()

		// routes[path](rec, req)
		router.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
})
