package utils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seadiaz/katalog/src/utils"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils")
}

var _ = Describe("Utils | Serialize", func() {
	It("should encode a string", func() {
		input := "dummy text"

		output := utils.Serialize(input)

		Expect(output).To(Equal("\"dummy text\""))
	})

	It("should encode a object", func() {
		type DummyObject struct {
			Key string
		}
		input := DummyObject{Key: "dummy text"}

		output := utils.Serialize(input)

		Expect(output).To(Equal("{\"Key\":\"dummy text\"}"))
	})

	It("shlud return empty string on encoding error", func() {
		input := func() {}

		output := utils.Serialize(input)

		Expect(output).To(Equal(""))
	})
})
