package publishers_test

import (
	"reflect"
	"testing"

	"github.com/avast/retry-go"
	"github.com/maxcnunes/httpfake"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/walmartdigital/katalog/src/collector/publishers"
	"github.com/walmartdigital/katalog/src/domain"
)

func retryDoDouble(retryableFunc retry.RetryableFunc, opts ...retry.Option) error {
	return retryableFunc()
}

// DummyK8sResource ...
type DummyK8sResource struct {
	ID string `json:",omitempty"`
}

// GetID ...
func (s *DummyK8sResource) GetID() string {
	return s.ID
}

// GetType ...
func (s *DummyK8sResource) GetType() reflect.Type {
	return reflect.TypeOf(s)
}

// GetK8sResource ...
func (s *DummyK8sResource) GetK8sResource() interface{} {
	return s
}

// GetGeneration ...
func (s *DummyK8sResource) GetGeneration() int64 {
	return s.GetGeneration()
}

// GetNamespace ...
func (s *DummyK8sResource) GetNamespace() string {
	return s.GetNamespace()
}

// GetName ...
func (s *DummyK8sResource) GetName() string {
	return s.GetName()
}

// GetAnnotations ...
func (s *DummyK8sResource) GetAnnotations() map[string]string {
	return s.GetAnnotations()
}

// GetLabels ...
func (s *DummyK8sResource) GetLabels() map[string]string {
	return s.GetLabels()
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

func createCreateFakeServer(path string, statusCode int, body string) *httpfake.HTTPFake {
	output := httpfake.New()
	output.
		NewHandler().
		Post(path).
		Reply(statusCode).
		BodyString(body)
	return output
}

var _ = Describe("add", func() {
	It("should return nil error when add service request succeed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		body := `{"ID": "` + serviceID + `"}`
		fakeService := createCreateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.Service{ID: serviceID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when add service request fail with status code 404", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 404
		body := `{"status": "fail"}`
		fakeService := createCreateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.Service{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("post service failed"))
	})

	It("should return an error when add service request fail with status code 500", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeService := createCreateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.Service{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("post service failed"))
	})

	It("should return nil error when add deployment request succeed", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 200
		body := `{"ID": "` + deploymentID + `"}`
		fakeDeployment := createCreateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{ID: deploymentID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when add deployment request fail with status code 404", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 404
		body := `{"status": "fail"}`
		fakeDeployment := createCreateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("post deployment failed"))
	})

	It("should return an error when add deployment request fail with status code 500", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeDeployment := createCreateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("post deployment failed"))
	})

	It("should return nil error when add statefulset request succeed", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 200
		body := `{"ID": "` + statefulSetID + `"}`
		fakeStatefulSet := createCreateFakeServer(path, statusCode, body)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{ID: statefulSetID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when add statefulset request fail with status code 404", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 404
		body := `{"status": "fail"}`
		fakeStatefulSet := createCreateFakeServer(path, statusCode, body)
		defer fakeStatefulSet.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("post statefulset failed"))
	})

	It("should return an error when add statefulset request fail with status code 500", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeStatefulSet := createCreateFakeServer(path, statusCode, body)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("post statefulset failed"))
	})

	It("should return nil if resource Type is not handled", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeService := createCreateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind:     domain.OperationTypeAdd,
			Resource: domain.Resource{K8sResource: new(DummyK8sResource)},
		})

		Expect(output).To(BeNil())
	})
})

func createUpdateFakeServer(path string, statusCode int, body string) *httpfake.HTTPFake {
	output := httpfake.New()
	output.
		NewHandler().
		Put(path).
		Reply(statusCode).
		BodyString(body)
	return output
}

var _ = Describe("update", func() {
	It("should return nil error when update service request succeed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		body := `{"ID": "` + serviceID + `"}`
		fakeService := createUpdateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.Service{ID: serviceID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when update service request fail with status code 404", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 404
		body := `{"status": "fail"}`
		fakeService := createUpdateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.Service{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put service failed"))
	})

	It("should return an error when update service request fail with status code 500", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeService := createUpdateFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.Service{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put service failed"))
	})

	It("should return nil error when update deployment request succeed", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 200
		body := `{"ID": "` + deploymentID + `"}`
		fakeDeployment := createUpdateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{ID: deploymentID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when update deployment request fail with status code 404", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 404
		body := `{"status": "fail"}`
		fakeDeployment := createUpdateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put deployment failed"))
	})

	It("should return an error when update deployment request fail with status code 500", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeDeployment := createUpdateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put deployment failed"))
	})

	It("should return nil error when update statefulset request succeed", func() {
		statefulSettID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSettID
		statusCode := 200
		body := `{"ID": "` + statefulSettID + `"}`
		fakeStatefulSet := createUpdateFakeServer(path, statusCode, body)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{ID: statefulSettID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when update statefulset request fail with status code 404", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + deploymentID
		statusCode := 404
		body := `{"status": "fail"}`
		fakeStatefulSet := createUpdateFakeServer(path, statusCode, body)
		defer fakeStatefulSet.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put statefulset failed"))
	})

	It("should return an error when update statefulset request fail with status code 500", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeStatefulSet := createUpdateFakeServer(path, statusCode, body)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put statefulset failed"))
	})

	It("should return nil if resource Type is not handled", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeDeployment := createUpdateFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind:     domain.OperationTypeUpdate,
			Resource: domain.Resource{K8sResource: new(DummyK8sResource)},
		})

		Expect(output).To(BeNil())
	})
})

func createDeleteFakeServer(path string, statusCode int) *httpfake.HTTPFake {
	output := httpfake.New()
	output.
		NewHandler().
		Delete(path).
		Reply(statusCode).
		BodyString("")
	return output
}

var _ = Describe("delete", func() {
	It("should return nil error when delete service request succeed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.Service{ID: serviceID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when service request fail with status code 404", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 404
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.Service{ID: serviceID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete service failed"))
	})

	It("should return an error when service request fail with status code 500", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.Service{ID: serviceID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete service failed"))
	})

	It("should return nil error when deployment request succeed", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 200
		fakeDeployment := createDeleteFakeServer(path, statusCode)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{ID: deploymentID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when deployment request fail with status code 404", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 404
		fakeDeployment := createDeleteFakeServer(path, statusCode)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{ID: deploymentID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete deployment failed"))
	})

	It("should return an error when deployment request fail with status code 500", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 500
		fakeDeployment := createDeleteFakeServer(path, statusCode)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.Deployment{ID: deploymentID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete deployment failed"))
	})

	It("should return nil error when statefulset request succeed", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 200
		fakeStatefulSet := createDeleteFakeServer(path, statusCode)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{ID: statefulSetID},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when statefulset request fail with status code 404", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 404
		fakeStatefulSet := createDeleteFakeServer(path, statusCode)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{ID: statefulSetID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete statefulset failed"))
	})

	It("should return an error when statefulset request fail with status code 500", func() {
		statefulSetID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/statefulsets/" + statefulSetID
		statusCode := 500
		fakeStatefulSet := createDeleteFakeServer(path, statusCode)
		defer fakeStatefulSet.Server.Close()
		url := fakeStatefulSet.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &domain.StatefulSet{ID: statefulSetID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete statefulset failed"))
	})

	It("should return nil if resource Type is not handled", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Resource: domain.Resource{
				K8sResource: &DummyK8sResource{ID: serviceID},
			},
		})

		Expect(output).To(BeNil())
	})
})

var _ = Describe("unknown", func() {
	It("should return error when operation is unknown", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: "unknown",
			Resource: domain.Resource{
				K8sResource: &domain.Service{ID: serviceID},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("operation unknown"))
	})
})
