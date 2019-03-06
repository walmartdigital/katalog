package persistence

// Persistence ...
type Persistence interface {
	Create(kind string, id string, obj interface{})
	GetAll(kind string) []interface{}
	Close()
}
