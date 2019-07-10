package repositories

// Repository ...
type Repository interface {
	CreateResource(obj interface{}) error
	UpdateResource(obj interface{}) error
	DeleteResource(obj interface{}) error
	GetAllResources() []interface{}
}
