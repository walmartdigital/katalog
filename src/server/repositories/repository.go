package repositories

// Repository ...
type Repository interface {
	CreateService(obj interface{})
	DeleteService(obj interface{})
	GetAllServices() []interface{}
}
