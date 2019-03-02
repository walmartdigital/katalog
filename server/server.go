package server

import (
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

// Article ...
type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Articles ...
type Articles []Article

// Run ...
func Run() {
	glog.Info("server starging...")
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/services", getAllServices).Methods("GET")
	myRouter.HandleFunc("/services/{id}", createService).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
