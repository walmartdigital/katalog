package persistence

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DummyMemory struct{}

func (m *DummyMemory) Add(values ...interface{}) {
	fmt.Println("oh yeah!")
}

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "memory persistence")
}

var _ = Describe("create", func() {
	It("should add a single struct", func() {
		memory := make(map[string]interface{})
		persistence := BuildMemoryPersistence(memory)
		value := struct{ name string }{"deadpool"}

		persistence.Create("dummy", "myid", value)

		Expect(memory["dummy-myid"]).To(Equal(value))
	})
})

var _ = Describe("get all", func() {
	It("should return all values", func() {
		memory := make(map[string]interface{})
		persistence := BuildMemoryPersistence(memory)
		value := struct{ name string }{"deadpool"}
		memory["dummy-myid"] = value
		expected := make([]interface{}, 1)
		expected[0] = value

		results := persistence.GetAll("dummy")

		Expect(results).To(Equal(expected))
	})
})

var _ = Describe("delete", func() {
	It("should delete an item", func() {
		memory := make(map[string]interface{})
		persistence := BuildMemoryPersistence(memory)
		value := struct{ name string }{"deadpool"}
		memory["dummy-myid"] = value

		persistence.Delete("dummy", "myid")

		Expect(memory["dummy-myid"]).To(BeNil())
	})
})
