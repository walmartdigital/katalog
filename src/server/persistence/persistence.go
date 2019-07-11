package persistence

// Persistence ...
type Persistence interface {
	Get(id string) (interface{}, error)
	Create(id string, obj interface{}) error
	Update(id string, obj interface{}) error
	Delete(id string) error
	GetAll() []interface{}
}
