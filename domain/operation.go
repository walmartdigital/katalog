package domain

const (
	// OperationTypeAdd ...
	OperationTypeAdd OperationType = "add"
	// OperationTypeDelete ...
	OperationTypeDelete OperationType = "delete"
	// OperationTypeUpdate ...
	OperationTypeUpdate OperationType = "update"
)

// OperationType ...
type OperationType string

// Operation ...
type Operation struct {
	Kind    OperationType `json:",omitempty"`
	Service Service
}
