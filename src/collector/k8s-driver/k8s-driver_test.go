package k8sdriver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8sdriver "github.com/walmartdigital/katalog/src/collector/k8s-driver"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

var _ = Describe("build driver", func() {
	It("should build a driver", func() {
		kubePath := ""
		excludeSystemNamespace := true

		output := k8sdriver.BuildDriver(kubePath, excludeSystemNamespace)

		Expect(output)
	})

})
