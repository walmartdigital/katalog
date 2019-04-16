package persistence

// Persistence ...
type Persistence interface {
	Create(kind string, id string, obj interface{})
	Delete(kind string, id string)
	GetAll(kind string) []interface{}
}
