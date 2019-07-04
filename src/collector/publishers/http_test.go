package publishers_test

import (
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

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

var _ = Describe("create", func() {
	It("should return nil error when service request succeed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		body := `{"ID": "` + serviceID + `"}`
		fakeService := createFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when service request failed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeService := createFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put service failed"))
	})

	It("should return an error when service request respond 500", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeService := createFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put service failed"))
	})

	It("should return nil error when deployment request succeed", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 200
		body := `{"ID": "` + deploymentID + `"}`
		fakeDeployment := createFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when deployment request failed", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeDeployment := createFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := "localhost:5000"
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put deployment failed"))
	})

	It("should return an error when deployment request respond 500", func() {
		deploymentID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/deployments/" + deploymentID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeDeployment := createFakeServer(path, statusCode, body)
		defer fakeDeployment.Server.Close()
		url := fakeDeployment.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put deployment failed"))
	})

	It("should return nil if resource Type is not handled", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		body := `{"status": "fail"}`
		fakeService := createFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Resource: domain.Resource{
				Type:   "XXXX",
				Object: domain.Service{},
			},
		})

		Expect(output).To(BeNil())
	})
})

func createFakeServer(path string, statusCode int, body string) *httpfake.HTTPFake {
	output := httpfake.New()
	output.
		NewHandler().
		Put(path).
		Reply(statusCode).
		BodyString(body)
	return output
}

var _ = Describe("delete", func() {
	It("should return nil error when service request succeed", func() {
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
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
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
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
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
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
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
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
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
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
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
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
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
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
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
				Type: "Deployment",
				Object: domain.Deployment{
					ID: deploymentID,
				},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete deployment failed"))
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
				Type: "XXXX",
				Object: domain.Service{
					ID: serviceID,
				},
			},
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

var _ = Describe("update", func() {
	It("should return nil error when request succeed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		body := `{"ID": "` + serviceID + `"}`
		fakeService := createFakeServer(path, statusCode, body)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeUpdate,
			Resource: domain.Resource{
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
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
				Type: "Service",
				Object: domain.Service{
					ID: serviceID,
				},
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("operation unknown"))
	})
})
