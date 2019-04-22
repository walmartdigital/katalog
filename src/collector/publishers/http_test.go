package publishers_test

import (
	"testing"

	"github.com/maxcnunes/httpfake"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/walmartdigital/katalog/src/collector/publishers"
	"github.com/walmartdigital/katalog/src/domain"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

var _ = Describe("create", func() {
	It("should add a service", func() {
		fakeService := httpfake.New()
		defer fakeService.Server.Close()
		serviceID := "6425377e-badd-4c46-828a-00c9afa7a156"
		fakeService.NewHandler().
			Put("/services/" + serviceID).
			Reply(200).
			BodyString(`{"ID": "` + serviceID + `"}`)
		url := fakeService.ResolveURL("")
		publisher := publishers.BuildHTTPPublisher(url)

		output := publisher.Publish(domain.Operation{
			Kind: domain.OperationTypeAdd,
			Service: domain.Service{
				ID: serviceID,
			},
		})

		Expect(output).To(BeNil())
	})
})
