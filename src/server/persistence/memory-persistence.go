package persistence

import "github.com/emirpasic/gods/lists/arraylist"

// MemoryPersistence is a memory implementantion of persistence
type MemoryPersistence struct {
	memory map[string]interface{}
}

// BuildMemoryPersistence ...
func BuildMemoryPersistence(memory map[string]interface{}) Persistence {
	return &MemoryPersistence{
		memory: memory,
	}
}

// Create ...
func (p *MemoryPersistence) Create(kind string, id string, obj interface{}) {
	p.memory[kind+"-"+id] = obj
}

// Delete ...
func (p *MemoryPersistence) Delete(kind string, id string) {
	delete(p.memory, kind+"-"+id)
}

// GetAll ...
func (p *MemoryPersistence) GetAll(kind string) []interface{} {
	list := arraylist.New()
	for _, value := range p.memory {
		list.Add(value)
	}
	return list.Values()
}
