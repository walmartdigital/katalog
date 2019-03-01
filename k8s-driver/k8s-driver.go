package k8sdriver

import (
	"github.com/golang/glog"
	"github.com/seadiaz/katalog/domain"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const resyncPeriod = 0

// Driver ...
type Driver struct {
	clientSet              *kubernetes.Clientset
	excludeSysmteNamespace bool
}

// BuildDriver ...
func BuildDriver(kubeconfigPath string, excludeSysmteNamespace bool) *Driver {
	return &Driver{
		clientSet:              buildClientSet(kubeconfigPath),
		excludeSysmteNamespace: excludeSysmteNamespace,
	}
}

// StartWatchingServices ...
func (d *Driver) StartWatchingServices(events chan interface{}) {
	watchList := d.buildWatchListForServices(v1.ResourceServices)
	controller := d.buildController(watchList, d.createAddHandler(events), d.createUpdateHandler(events), d.createDeleteHandler(events))
	controller.Run(make(chan struct{}))
}

func buildClientSet(kubeconfigPath string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		glog.Errorln(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorln(err)
	}
	return clientset
}

func (d *Driver) buildWatchListForServices(resource v1.ResourceName) *cache.ListWatch {
	watchlist := cache.NewListWatchFromClient(
		d.clientSet.CoreV1().RESTClient(),
		string(resource),
		v1.NamespaceAll,
		fields.Everything(),
	)
	return watchlist
}

func (d *Driver) buildController(watchList *cache.ListWatch, addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{})) cache.Controller {
	_, controller := cache.NewInformer(
		watchList,
		&v1.Service{},
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			UpdateFunc: updateFunc,
			DeleteFunc: deleteFunc,
		},
	)
	return controller
}

func (d *Driver) createDeleteHandler(channel chan interface{}) func(interface{}) {
	return func(obj interface{}) {
		k8sService := obj.(*v1.Service)
		endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
		service := buildOperationFromK8sService(domain.OperationTypeDelete, k8sService, *endpoints)
		channel <- service
	}
}

func (d *Driver) createAddHandler(channel chan interface{}) func(interface{}) {
	return func(obj interface{}) {
		k8sService := obj.(*v1.Service)
		if d.excludeSysmteNamespace && k8sService.Namespace == "kube-system" {
			glog.Infof("%s excluded because belongs to kube-system namespace", k8sService.Name)
			return
		}
		endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
		service := buildOperationFromK8sService(domain.OperationTypeAdd, k8sService, *endpoints)
		channel <- service
	}
}

func (d *Driver) createUpdateHandler(channel chan interface{}) func(oldObj interface{}, newObj interface{}) {
	return func(oldObj interface{}, newObj interface{}) {
		k8sService := newObj.(*v1.Service)
		endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
		service := buildOperationFromK8sService(domain.OperationTypeUpdate, k8sService, *endpoints)
		channel <- service
	}
}
