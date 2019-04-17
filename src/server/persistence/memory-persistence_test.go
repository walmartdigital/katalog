package persistence_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walmartdigital/katalog/src/server/persistence"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

var _ = Describe("create", func() {
	It("should add a single struct", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		value := struct{ name string }{"deadpool"}

		persistence.Create("dummy", "myid", value)

		Expect(memory["dummy-myid"]).To(Equal(value))
	})
})

var _ = Describe("get all", func() {
	It("should return all values", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		value := struct{ name string }{"deadpool"}
		memory["dummy-myid"] = value
		expected := make([]interface{}, 1)
		expected[0] = value

		results := persistence.GetAll("dummy")

		Expect(results).To(Equal(expected))
	})
})

var _ = Describe("delete", func() {
	It("should remove one item", func() {
		memory := make(map[string]interface{})
		persistence := persistence.BuildMemoryPersistence(memory)
		value := struct{ name string }{"max"}
		memory["dummy-myid"] = value
		expected := make([]interface{}, 1)
		expected[0] = value

		persistence.Delete("dummy", "myid")

		Expect(len(memory)).To(Equal(0))
		Expect(memory["dummy-myid"]).To(BeNil())
	})
})
