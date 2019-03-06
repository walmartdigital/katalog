package repositories

// Repository ...
type Repository interface {
	CreateService(obj interface{})
	GetAllServices() []interface{}
}
