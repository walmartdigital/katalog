package server

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var mutex = &sync.Mutex{}
var metrics map[string]interface{}

// Metrics ...
type Metrics interface {
	InitMetrics()
	IncrementCounter(string, ...string)
	DestroyMetrics()
}

// MetricsFactory ...
type MetricsFactory interface {
	Create() Metrics
}

// PrometheusMetrics ...
type PrometheusMetrics struct {
}

// InitMetrics ...
func (p PrometheusMetrics) InitMetrics() {
	mutex.Lock()
	if metrics == nil {
		metrics = make(map[string]interface{})

		metrics["createDeployment"] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "katalog",
				Subsystem: "deployment",
				Name:      "create",
				Help:      "Total number of deployment creations",
			},
			[]string{"id", "ns", "rn"},
		)
		metrics["createDeployment"].(*prometheus.CounterVec).WithLabelValues("", "", "")
		prometheus.MustRegister(metrics["createDeployment"].(*prometheus.CounterVec))

		metrics["updateDeployment"] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "katalog",
				Subsystem: "deployment",
				Name:      "update",
				Help:      "Total number of deployment updates",
			},
			[]string{"id", "ns", "rn"},
		)
		metrics["updateDeployment"].(*prometheus.CounterVec).WithLabelValues("", "", "")
		prometheus.MustRegister(metrics["updateDeployment"].(*prometheus.CounterVec))

		metrics["deleteDeployment"] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "katalog",
				Subsystem: "deployment",
				Name:      "delete",
				Help:      "Total number of deployment deletes",
			},
			[]string{"id", "ns", "rn"},
		)
		metrics["deleteDeployment"].(*prometheus.CounterVec).WithLabelValues("", "", "")
		prometheus.MustRegister(metrics["deleteDeployment"].(*prometheus.CounterVec))

		metrics["createStatefulSet"] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "katalog",
				Subsystem: "statefulset",
				Name:      "create",
				Help:      "Total number of statefulset creations",
			},
			[]string{"id", "ns", "rn"},
		)
		metrics["createStatefulSet"].(*prometheus.CounterVec).WithLabelValues("", "", "")
		prometheus.MustRegister(metrics["createStatefulSet"].(*prometheus.CounterVec))

		metrics["updateStatefulSet"] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "katalog",
				Subsystem: "statefulset",
				Name:      "update",
				Help:      "Total number of statefulset updates",
			},
			[]string{"id", "ns", "rn"},
		)
		metrics["updateStatefulSet"].(*prometheus.CounterVec).WithLabelValues("", "", "")
		prometheus.MustRegister(metrics["updateStatefulSet"].(*prometheus.CounterVec))

		metrics["deleteStatefulSet"] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "katalog",
				Subsystem: "statefulset",
				Name:      "delete",
				Help:      "Total number of statefulset deletes",
			},
			[]string{"id", "ns", "rn"},
		)
		metrics["deleteStatefulSet"].(*prometheus.CounterVec).WithLabelValues("", "", "")
		prometheus.MustRegister(metrics["deleteStatefulSet"].(*prometheus.CounterVec))
	}
	mutex.Unlock()
}

// IncrementCounter ...
func (p PrometheusMetrics) IncrementCounter(key string, labels ...string) {
	metrics[key].(*prometheus.CounterVec).WithLabelValues(labels...).Inc()
}

// DestroyMetrics ...
func (p PrometheusMetrics) DestroyMetrics() {
	mutex.Lock()
	for _, v := range metrics {
		prometheus.Unregister(v.(prometheus.Collector))
	}
	mutex.Unlock()
}
