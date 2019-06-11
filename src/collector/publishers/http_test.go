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
			Kind: domain.OperationTypeAdd,
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when request failed", func() {
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
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put service failed"))
	})

	It("should return an error when request respond 500", func() {
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
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("put service failed"))
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
	It("should return nil error when request succeed", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 200
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).To(BeNil())
	})

	It("should return an error when request fail with status code 404", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 404
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete service failed"))
	})

	It("should return an error when request fail with status code 500", func() {
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		path := "/services/" + serviceID
		statusCode := 500
		fakeService := createDeleteFakeServer(path, statusCode)
		defer fakeService.Server.Close()
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url, retryDoDouble)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeDelete,
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("delete service failed"))
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
			Service: domain.Service{
				ID: serviceID,
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
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).ToNot(BeNil())
		Expect(output.Error()).To(Equal("operation unknown"))
	})
})
