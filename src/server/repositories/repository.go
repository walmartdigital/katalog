package repositories

// Repository ...
type Repository interface {
	CreateService(obj interface{}) error
	DeleteService(obj interface{}) error
	GetAllServices() []interface{}
}
