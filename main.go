package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/seadiaz/katalog/k8s-driver"
	"github.com/seadiaz/katalog/publishers"
)

var consulAddress = flag.String("consul-addr", "127.0.0.1:8500", "consul address")

func main() {
	flag.Parse()
	glog.Info("[main] starging...")

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	mainCollector(kubeconfig)
}

func mainCollector(kubeconfig string) {
	serviceEvents := make(chan interface{})
	k8sDriver := k8sdriver.BuildDriver(kubeconfig)
	publisher := publishers.Create(*consulAddress)
	go k8sDriver.StartWatchingServices(serviceEvents)
	for {
		select {
		case event := <-serviceEvents:
			publisher.Publish(event)
		}
	}
}
