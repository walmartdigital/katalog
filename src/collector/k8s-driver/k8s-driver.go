package k8sdriver

import (
	"reflect"

	"github.com/golang/glog"
	"github.com/walmartdigital/katalog/src/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	excludeSystemNamespace bool
}

// BuildDriver ...
func BuildDriver(kubeconfigPath string, excludeSystemNamespace bool) *Driver {
	return &Driver{
		clientSet:              buildClientSet(kubeconfigPath),
		excludeSystemNamespace: excludeSystemNamespace,
	}
}

// StartWatchingResources ...
func (d *Driver) StartWatchingResources(events chan interface{}, resource domain.Resource) {
	listWatch := d.buildListWatchForResources(resource)
	controller := d.buildController(listWatch, resource, d.createAddHandler(events, resource), d.createUpdateHandler(events, resource), d.createDeleteHandler(events, resource))
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

func (d *Driver) buildListWatchForResources(resource domain.Resource) *cache.ListWatch {
	var listWatch *cache.ListWatch

	switch v := resource.GetType(); v {
	case reflect.TypeOf(domain.Service{}):
		listWatch = cache.NewListWatchFromClient(
			d.clientSet.CoreV1().RESTClient(),
			"services",
			corev1.NamespaceAll,
			fields.Everything(),
		)
	case reflect.TypeOf(domain.Deployment{}):
		listWatch = cache.NewListWatchFromClient(
			d.clientSet.AppsV1().RESTClient(),
			"deployments",
			corev1.NamespaceAll,
			fields.Everything(),
		)
	default:
		glog.Errorf("Type %s not found", v)
	}

	return listWatch
}

func (d *Driver) buildController(listWatch *cache.ListWatch, resource domain.Resource, addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{})) cache.Controller {
	var controller cache.Controller

	switch v := resource.GetType(); v {
	case reflect.TypeOf(domain.Service{}):
		_, controller = cache.NewInformer(
			listWatch,
			&corev1.Service{},
			resyncPeriod,
			cache.ResourceEventHandlerFuncs{
				AddFunc:    addFunc,
				UpdateFunc: updateFunc,
				DeleteFunc: deleteFunc,
			},
		)
	case reflect.TypeOf(domain.Deployment{}):
		_, controller = cache.NewInformer(
			listWatch,
			&appsv1.Deployment{},
			resyncPeriod,
			cache.ResourceEventHandlerFuncs{
				AddFunc:    addFunc,
				UpdateFunc: updateFunc,
				DeleteFunc: deleteFunc,
			},
		)
	default:
		glog.Errorf("Type %s not found", v)
	}

	return controller
}

func (d *Driver) createAddHandler(channel chan interface{}, resource domain.Resource) func(interface{}) {
	return func(obj interface{}) {
		if resource.GetType() == reflect.TypeOf(domain.Service{}) {
			k8sService := obj.(*corev1.Service)
			if d.excludeSystemNamespace && k8sService.Namespace == "kube-system" {
				glog.Infof("%s excluded because belongs to kube-system namespace", k8sService.Name)
				return
			}
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeAdd, k8sService, *endpoints)
			channel <- service
		}

		if resource.GetType() == reflect.TypeOf(domain.Deployment{}) {
			k8sDeployment := obj.(*appsv1.Deployment)
			if d.excludeSystemNamespace && k8sDeployment.Namespace == "kube-system" {
				glog.Infof("%s excluded because belongs to kube-system namespace", k8sDeployment.Name)
				return
			}
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeAdd, k8sDeployment)
			channel <- deployment
		}
	}
}

func (d *Driver) createDeleteHandler(channel chan interface{}, resource domain.Resource) func(interface{}) {
	return func(obj interface{}) {
		if resource.GetType() == reflect.TypeOf(domain.Service{}) {
			k8sService := obj.(*corev1.Service)
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeDelete, k8sService, *endpoints)
			channel <- service
		}
		if resource.GetType() == reflect.TypeOf(domain.Deployment{}) {
			k8sDeployment := obj.(*appsv1.Deployment)
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeDelete, k8sDeployment)
			channel <- deployment
		}
	}
}

func (d *Driver) createUpdateHandler(channel chan interface{}, resource domain.Resource) func(oldObj interface{}, newObj interface{}) {
	return func(oldObj interface{}, newObj interface{}) {
		if resource.GetType() == reflect.TypeOf(domain.Service{}) {
			k8sService := newObj.(*corev1.Service)
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeUpdate, k8sService, *endpoints)
			channel <- service
		}
		if resource.GetType() == reflect.TypeOf(domain.Deployment{}) {
			k8sDeployment := newObj.(*appsv1.Deployment)
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeUpdate, k8sDeployment)
			channel <- deployment
		}
	}
}
