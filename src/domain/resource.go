package domain

const (
	// ResourceKTypeService ...
	ResourceKTypeService ResourceType = "Service"
	// ResourceKTypeDeployment ...
	ResourceKTypeDeployment ResourceType = "Deployment"
)

// ResourceType ...
type ResourceType string

// Resource ...
type Resource struct {
	Type   ResourceType `json:"type"`
	Object interface{}
}
