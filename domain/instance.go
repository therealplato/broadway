package domain

import (
	"encoding/json"

	"github.com/namely/broadway/store"
)

// Instance entity
type Instance struct {
	PlaybookID string            `json:"playbook_id" binding:"required"`
	ID         string            `json:"id"`
	Created    string            `json:"created"`
	Vars       map[string]string `json:"vars"`
}

// Instance json representation
func (i *Instance) JSON() (string, error) {
	encoded, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// Path for an instance
func (i *Instance) Path() string {
	return "/broadway/instances/" + i.PlaybookID + "/" + i.ID
}

// Instance repository interface
type InstanceRepository interface {
	Save(instance Instance) error
}

// Instance repository to handle persistence logic
type InstanceRepo struct {
	store store.Store
}

// Create a new instance repo
func NewInstanceRepo(s store.Store) *InstanceRepo {
	return &InstanceRepo{store: s}
}

// Save a new instance as json
func (ir *InstanceRepo) Save(instance Instance) error {
	encoded, err := instance.JSON()
	if err != nil {
		return err
	}
	err = ir.store.SetValue(instance.Path(), encoded)
	if err != nil {
		return err
	}
	return nil
}
