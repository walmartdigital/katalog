package k8sdriver

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/walmartdigital/katalog/src/domain"
	"github.com/walmartdigital/katalog/src/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var log = logrus.New()

func init() {
	err := utils.LogInit(log)
	if err != nil {
		logrus.Fatal(err)
	}
}

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
		log.Errorln(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorln(err)
	}
	return clientset
}

func (d *Driver) buildListWatchForResources(resource domain.Resource) *cache.ListWatch {
	var listWatch *cache.ListWatch

	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
		listWatch = cache.NewListWatchFromClient(
			d.clientSet.CoreV1().RESTClient(),
			"services",
			corev1.NamespaceAll,
			fields.Everything(),
		)
	case reflect.TypeOf(new(domain.Deployment)):
		listWatch = cache.NewListWatchFromClient(
			d.clientSet.AppsV1().RESTClient(),
			"deployments",
			corev1.NamespaceAll,
			fields.Everything(),
		)
	case reflect.TypeOf(new(domain.StatefulSet)):
		listWatch = cache.NewListWatchFromClient(
			d.clientSet.AppsV1().RESTClient(),
			"statefulsets",
			corev1.NamespaceAll,
			fields.Everything(),
		)
	default:
		log.Errorf("Type %s not found", v)
	}

	return listWatch
}

func (d *Driver) buildController(listWatch *cache.ListWatch, resource domain.Resource, addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{})) cache.Controller {
	var controller cache.Controller

	switch v := resource.GetType(); v {
	case reflect.TypeOf(new(domain.Service)):
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
	case reflect.TypeOf(new(domain.Deployment)):
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
	case reflect.TypeOf(new(domain.StatefulSet)):
		_, controller = cache.NewInformer(
			listWatch,
			&appsv1.StatefulSet{},
			resyncPeriod,
			cache.ResourceEventHandlerFuncs{
				AddFunc:    addFunc,
				UpdateFunc: updateFunc,
				DeleteFunc: deleteFunc,
			},
		)
	default:
		log.Errorf("Type %s not found", v)
	}

	return controller
}

func (d *Driver) createAddHandler(channel chan interface{}, resource domain.Resource) func(interface{}) {
	return func(obj interface{}) {
		switch t := resource.GetType(); t {
		case reflect.TypeOf(new(domain.Service)):
			k8sService := obj.(*corev1.Service)
			if d.excludeSystemNamespace && k8sService.Namespace == "kube-system" {
				log.Infof("%s excluded because belongs to kube-system namespace", k8sService.Name)
				return
			}
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeAdd, k8sService, *endpoints)
			channel <- service
		case reflect.TypeOf(new(domain.Deployment)):
			k8sDeployment := obj.(*appsv1.Deployment)
			if d.excludeSystemNamespace && k8sDeployment.Namespace == "kube-system" {
				log.Infof("%s excluded because belongs to kube-system namespace", k8sDeployment.Name)
				return
			}
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeAdd, k8sDeployment)
			channel <- deployment
		case reflect.TypeOf(new(domain.StatefulSet)):
			k8sStatefulSet := obj.(*appsv1.StatefulSet)
			if d.excludeSystemNamespace && k8sStatefulSet.Namespace == "kube-system" {
				log.Infof("%s excluded because belongs to kube-system namespace", k8sStatefulSet.Name)
				return
			}
			statefulset := buildOperationFromK8sStatefulSet(domain.OperationTypeAdd, k8sStatefulSet)
			channel <- statefulset
		default:
			log.Errorf("Type %s not found", t)
		}
	}
}

func (d *Driver) createDeleteHandler(channel chan interface{}, resource domain.Resource) func(interface{}) {
	return func(obj interface{}) {
		switch t := resource.GetType(); t {
		case reflect.TypeOf(new(domain.Service)):
			k8sService := obj.(*corev1.Service)
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeDelete, k8sService, *endpoints)
			channel <- service
		case reflect.TypeOf(new(domain.Deployment)):
			k8sDeployment := obj.(*appsv1.Deployment)
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeDelete, k8sDeployment)
			channel <- deployment
		case reflect.TypeOf(new(domain.StatefulSet)):
			k8sStatefulSet := obj.(*appsv1.StatefulSet)
			statefulset := buildOperationFromK8sStatefulSet(domain.OperationTypeDelete, k8sStatefulSet)
			channel <- statefulset
		default:
			log.Errorf("Type %s not found", t)
		}
	}
}

func (d *Driver) createUpdateHandler(channel chan interface{}, resource domain.Resource) func(oldObj interface{}, newObj interface{}) {
	return func(oldObj interface{}, newObj interface{}) {
		switch t := resource.GetType(); t {
		case reflect.TypeOf(new(domain.Service)):
			k8sService := newObj.(*corev1.Service)
			endpoints, _ := d.clientSet.CoreV1().Endpoints(k8sService.Namespace).Get(k8sService.Name, metav1.GetOptions{})
			service := buildOperationFromK8sService(domain.OperationTypeUpdate, k8sService, *endpoints)
			channel <- service
		case reflect.TypeOf(new(domain.Deployment)):
			k8sDeployment := newObj.(*appsv1.Deployment)
			deployment := buildOperationFromK8sDeployment(domain.OperationTypeUpdate, k8sDeployment)
			channel <- deployment
		case reflect.TypeOf(new(domain.StatefulSet)):
			k8sStatefulSet := newObj.(*appsv1.StatefulSet)
			statefulset := buildOperationFromK8sStatefulSet(domain.OperationTypeUpdate, k8sStatefulSet)
			channel <- statefulset
		default:
			log.Errorf("Type %s not found", t)
		}
	}
}
