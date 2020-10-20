package utils_test

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/utils"
)

type DummyObject struct {
	Key string
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils")
}

var _ = Describe("Utils", func() {
	Describe("Serialize", func() {
		It("should encode a string", func() {
			input := "dummy text"

			output := utils.Serialize(input)

			Expect(output).To(Equal(`"dummy text"`))
		})

		It("should encode a object", func() {
			input := DummyObject{Key: "dummy text"}

			output := utils.Serialize(input)

			Expect(output).To(Equal(`{"Key":"dummy text"}`))
		})

		It("should return empty string on encoding error", func() {
			input := func() {}

			output := utils.Serialize(input)

			Expect(output).To(Equal(""))
		})
	})

	Describe("Deserialize", func() {
		It("should return empty string on encoding error", func() {
			input := `"dummy text"`

			output := utils.Deserialize(input)

			Expect(output).To(Equal(""))
		})

		It("should decode a struct", func() {
			input := `{"Key":"dummy text"}`

			output := utils.Deserialize(input)

			outputParsed := DummyObject{}
			mapstructure.Decode(output, &outputParsed)
			Expect(outputParsed.Key).To(Equal("dummy text"))
		})
	})
})
