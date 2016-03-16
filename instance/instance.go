package instance

import "encoding/json"

// InstanceStatus represents the lifecycle state of one instance
type InstanceStatus string

const (
	InstanceStatusNew       InstanceStatus = ""
	InstanceStatusDeploying                = "deploying"
	InstanceStatusDeployed                 = "deployed"
	InstanceStatusDeleting                 = "deleting"
	InstanceStatusError                    = "error"
)

// InstanceAttributes contains metadata about an instance
type InstanceAttributes struct {
	PlaybookID string            `json:"playbook_id" binding:"required"`
	Id         string            `json:"id"`
	Created    string            `json:"created"`
	Vars       map[string]string `json:"vars"`
	Status     InstanceStatus    `json:"status"`
}

// JSON serializes a set of instance attributes
func (attrs *InstanceAttributes) JSON() (string, error) {
	encoded, err := json.Marshal(attrs)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// Instance represents an instantiation of a Playbook. The same Playbook might
// be used multiple times, e.g. for two similar pull requests on the same repo.
type Instance interface {
	json.Marshaler
	PlaybookID() string
	ID() string
	Save() error
	Destroy() error

	Attributes() *InstanceAttributes
	Status() InstanceStatus
}
