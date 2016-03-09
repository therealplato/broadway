package instance

import "encoding/json"

type InstanceStatus string

const (
	InstanceStatusNew       InstanceStatus = ""
	InstanceStatusDeploying                = "deploying"
	InstanceStatusDeployed                 = "deployed"
	InstanceStatusDeleting                 = "deleting"
	InstanceStatusError                    = "error"
)

type InstanceAttributes struct {
	PlaybookId string            `json:"playbook_id"`
	Id         string            `json:"id"`
	Created    string            `json:"created"`
	Vars       map[string]string `json:"vars"`
	Status     InstanceStatus    `json:"status"`
}

func (attrs *InstanceAttributes) JSON() (string, error) {
	encoded, err := json.Marshal(attrs)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

type Instance interface {
	PlaybookID() string
	ID() string
	Save() error
	Destroy() error

	Attributes() *InstanceAttributes
	Status() InstanceStatus
}
