package persistence

import (
	"errors"
	"sync"

	"github.com/emirpasic/gods/lists/arraylist"
)

// MemoryPersistence is a memory implementantion of persistence
type MemoryPersistence struct {
	memory *sync.Map
}

// BuildMemoryPersistence ...
func BuildMemoryPersistence(memory *sync.Map) Persistence {
	return &MemoryPersistence{
		memory: memory,
	}
}

// Get ...
func (p *MemoryPersistence) Get(id string) (interface{}, error) {
	if id == "" {
		return nil, errors.New("you must provide an id")
	}

	value, _ := p.memory.Load(id)

	return value, nil
}

// Create ...
func (p *MemoryPersistence) Create(id string, obj interface{}) error {
	if id == "" {
		return errors.New("you must provide an id")
	}
	p.memory.Store(id, obj)

	return nil
}

// Update ...
func (p *MemoryPersistence) Update(id string, obj interface{}) error {
	if id == "" {
		return errors.New("you must provide an id")
	}

	p.memory.Store(id, obj)

	return nil
}

// Delete ...
func (p *MemoryPersistence) Delete(id string) error {
	if id == "" {
		return errors.New("you must provide an id")
	}

	p.memory.Delete(id)

	return nil
}

// GetAll ...
func (p *MemoryPersistence) GetAll() ([]interface{}, error) {
	list := arraylist.New()

	p.memory.Range(func(key interface{}, value interface{}) bool {
		list.Add(value)

		return true
	})

	return list.Values(), nil
}
