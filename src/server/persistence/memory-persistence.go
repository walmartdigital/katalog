package persistence

import (
	"errors"

	"github.com/emirpasic/gods/lists/arraylist"
)

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
func (p *MemoryPersistence) Create(id string, obj interface{}) error {
	if id == "" {
		return errors.New("you must provide an id")
	}
	p.memory[id] = obj
	return nil
}

// Delete ...
func (p *MemoryPersistence) Delete(id string) error {
	if id == "" {
		return errors.New("you must provide an id")
	}
	delete(p.memory, id)
	return nil
}

// GetAll ...
func (p *MemoryPersistence) GetAll() []interface{} {
	list := arraylist.New()
	for _, value := range p.memory {
		list.Add(value)
	}
	return list.Values()
}
