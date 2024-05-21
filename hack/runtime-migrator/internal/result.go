package internal

type StatusType string

const (
	StatusSuccess           StatusType = "Success"
	StatusError             StatusType = "Error"
	StatusAlreadyExists     StatusType = "AlreadyExists"
	StatusRuntimeIdNotFound StatusType = "RuntimeIdNotFound"
)

type MigrationResult struct {
	RuntimeId    string     `json:"runtimeId"`
	ShootName    string     `json:"shootName"`
	Status       StatusType `json:"status"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
	PathToCRYaml string     `json:"pathToCRYaml,omitempty"`
}
