package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/seadiaz/katalog/domain"
)

func getAllServices(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	fmt.Println("Endpoint Hit: getAllServices")

	json.NewEncoder(w).Encode(articles)
}

func createService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var service domain.Service
	json.NewDecoder(r.Body).Decode(&service)

	glog.Info(service)

	fmt.Fprintf(w, "Key: "+key)
}
