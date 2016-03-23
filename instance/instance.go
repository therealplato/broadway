package instance

import (
	"encoding/json"
)

// Status represents the lifecycle state of one instance
type Status string

const (
	// StatusNew represents a newly created instance
	StatusNew Status = ""
	// StatusDeploying represents an instance that has begun deployment
	StatusDeploying = "deploying"
	// StatusDeployed represents an instance that has been deployed successfully
	StatusDeployed = "deployed"
	// StatusDeleting represents an instance that has begun deltion
	StatusDeleting = "deleting"
	// StatusError represents an instance that broke
	StatusError = "error"
)

// Attributes contains metadata about an instance
type Attributes struct {
	PlaybookID string            `json:"playbook_id" binding:"required"`
	ID         string            `json:"id"`
	Created    string            `json:"created"`
	Vars       map[string]string `json:"vars"`
	Status     Status            `json:"status"`
}

// JSON serializes a set of instance attributes
func (attrs *Attributes) JSON() (string, error) {
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

	Attributes() *Attributes
	Status() Status
}
