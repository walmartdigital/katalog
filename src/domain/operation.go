package domain

const (
	// OperationTypeAdd ...
	OperationTypeAdd OperationType = "add"
	// OperationTypeDelete ...
	OperationTypeDelete OperationType = "delete"
	// OperationTypeUpdate ...
	OperationTypeUpdate     OperationType = "update"
	ResourceKTypeService    ResourceType  = "Service"
	ResourceKTypeDeployment ResourceType  = "Deployment"
)

// OperationType ...
type OperationType string
type ResourceType string

// Operation ...
type Operation struct {
	Kind     OperationType `json:"kind"`
	Resource Resource      `json:"resource"`
}

type Resource struct {
	Type   ResourceType `json:"type"`
	Object interface{}
}
