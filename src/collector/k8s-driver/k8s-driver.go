package k8sdriver

import (
	"reflect"

	"github.com/golang/glog"
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const resyncPeriod = 0

// Driver ...
type Driver struct {
	clientSet              *kubernetes.Clientset
	excludeSystemNamespace bool
}

// BuildDriver ...
func BuildDriver(kubeconfigPath string, excludeSystemNamespace bool) *Driver {
	return &Driver{
		clientSet:              buildClientSet(kubeconfigPath),
		excludeSystemNamespace: excludeSystemNamespace,
	}
}

// StartWatchingServices ...
func (d *Driver) StartWatchingServices(events chan interface{}, resourceType runtime.Object) {
	watchList := d.buildWatchListForServices(corev1.ResourceServices)
	controller := d.buildController(watchList, resourceType, d.createAddHandler(events, resourceType), d.createUpdateHandler(events, resourceType), d.createDeleteHandler(events, resourceType))
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

func (d *Driver) buildWatchListForServices(resource corev1.ResourceName) *cache.ListWatch {
	watchlist := cache.NewListWatchFromClient(
		d.clientSet.CoreV1().RESTClient(),
		string(resource),
		corev1.NamespaceAll,
		fields.Everything(),
	)
	return watchlist
}

func (d *Driver) buildController(watchList *cache.ListWatch, resourceType runtime.Object, addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{})) cache.Controller {
	_, controller := cache.NewInformer(
		watchList,
		resourceType,
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			UpdateFunc: updateFunc,
			DeleteFunc: deleteFunc,
		},
	)
	return controller
}

func (d *Driver) createAddHandler(channel chan interface{}, resourceType runtime.Object) func(interface{}) {
	return func(obj interface{}) {
		if reflect.DeepEqual(resourceType, &corev1.Service{}) {
			k8sService := obj.(*corev1.Service)
			if d.excludeSystemNamespace && k8sService.Namespace == "kube-system" {
				glog.Infof("%s excluded because belongs to kube-system namespace", k8sService.Name)
				return
			}
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeAdd, k8sService, *endpoints)
			channel <- service
		}

		if reflect.DeepEqual(resourceType, &appsv1.Deployment{}) {
			k8sDeployment := obj.(*appsv1.Deployment)
			if d.excludeSystemNamespace && k8sDeployment.Namespace == "kube-system" {
				glog.Infof("%s excluded because belongs to kube-system namespace", k8sDeployment.Name)
				return
			}
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeAdd, resourceType)
			channel <- deployment
		}
	}
}

func (d *Driver) createDeleteHandler(channel chan interface{}, resourceType runtime.Object) func(interface{}) {
	return func(obj interface{}) {
		if reflect.DeepEqual(resourceType, &corev1.Service{}) {
			k8sService := obj.(*corev1.Service)
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeDelete, k8sService, *endpoints)
			channel <- service
		}
	}
}

func (d *Driver) createUpdateHandler(channel chan interface{}, resourceType runtime.Object) func(oldObj interface{}, newObj interface{}) {
	return func(oldObj interface{}, newObj interface{}) {
		if reflect.DeepEqual(resourceType, &corev1.Service{}) {
			k8sService := newObj.(*corev1.Service)
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeUpdate, k8sService, *endpoints)
			channel <- service
		}
	}
}
