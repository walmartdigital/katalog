package publishers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
	It("should add a single struct", func() {
		handler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("**")
			io.WriteString(w, "<html><body>Hello World!</body></html>")
		}

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		publisher := publishers.BuildHTTPPublisher("http://localhost")

		publisher.Publish(domain.Operation{Kind: domain.OperationTypeAdd})
	})
})
