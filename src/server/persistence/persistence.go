package persistence

// Persistence ...
type Persistence interface {
	Create(id string, obj interface{}) error
	Delete(id string) error
	GetAll() []interface{}
}
