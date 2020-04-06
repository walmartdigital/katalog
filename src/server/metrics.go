package server

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var mutex = &sync.Mutex{}
var currentMetricsMap *map[string]interface{}

// InitMetrics ...
func InitMetrics() *map[string]interface{} {
	mutex.Lock()

	if currentMetricsMap != nil {
		return currentMetricsMap
	}

	metricsmap := make(map[string]interface{})

	metricsmap["createDeployment"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "deployment",
			Name:      "create",
			Help:      "Total number of deployment creations",
		},
		[]string{"id", "ns", "rn"},
	)
	metricsmap["createDeployment"].(*prometheus.CounterVec).WithLabelValues("", "", "")
	prometheus.MustRegister(metricsmap["createDeployment"].(*prometheus.CounterVec))

	metricsmap["updateDeployment"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "deployment",
			Name:      "update",
			Help:      "Total number of deployment updates",
		},
		[]string{"id", "ns", "rn"},
	)
	metricsmap["updateDeployment"].(*prometheus.CounterVec).WithLabelValues("", "", "")
	prometheus.MustRegister(metricsmap["updateDeployment"].(*prometheus.CounterVec))

	metricsmap["deleteDeployment"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "deployment",
			Name:      "delete",
			Help:      "Total number of deployment deletes",
		},
		[]string{"id", "ns", "rn"},
	)
	metricsmap["deleteDeployment"].(*prometheus.CounterVec).WithLabelValues("", "", "")
	prometheus.MustRegister(metricsmap["deleteDeployment"].(*prometheus.CounterVec))

	metricsmap["createStatefulSet"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "statefulset",
			Name:      "create",
			Help:      "Total number of statefulset creations",
		},
		[]string{"id", "ns", "rn"},
	)
	metricsmap["createStatefulSet"].(*prometheus.CounterVec).WithLabelValues("", "", "")
	prometheus.MustRegister(metricsmap["createStatefulSet"].(*prometheus.CounterVec))

	metricsmap["updateStatefulSet"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "statefulset",
			Name:      "update",
			Help:      "Total number of statefulset updates",
		},
		[]string{"id", "ns", "rn"},
	)
	metricsmap["updateStatefulSet"].(*prometheus.CounterVec).WithLabelValues("", "", "")
	prometheus.MustRegister(metricsmap["updateStatefulSet"].(*prometheus.CounterVec))

	metricsmap["deleteStatefulSet"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "katalog",
			Subsystem: "statefulset",
			Name:      "delete",
			Help:      "Total number of statefulset deletes",
		},
		[]string{"id", "ns", "rn"},
	)
	metricsmap["deleteStatefulSet"].(*prometheus.CounterVec).WithLabelValues("", "", "")
	prometheus.MustRegister(metricsmap["deleteStatefulSet"].(*prometheus.CounterVec))

	currentMetricsMap = &metricsmap

	mutex.Unlock()
	return &metricsmap
}

// DestroyMetrics ...
func DestroyMetrics() {
	for _, v := range *(currentMetricsMap) {
		prometheus.Unregister(v.(prometheus.Collector))
	}
}
