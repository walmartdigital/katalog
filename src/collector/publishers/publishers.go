package publishers

// Publisher ...
type Publisher interface {
	Publish(obj interface{}) error
}
