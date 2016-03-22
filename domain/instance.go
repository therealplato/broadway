package domain

import (
	"encoding/json"

	"github.com/namely/broadway/store"
)

type Instance struct {
	PlaybookID string            `json:"playbook_id" binding:"required"`
	ID         string            `json:"id"`
	Created    string            `json:"created"`
	Vars       map[string]string `json:"vars"`
}

func (i *Instance) JSON() (string, error) {
	encoded, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (i *Instance) Path() string {
	return "/broadway/instances/" + i.PlaybookID + "/" + i.ID
}

type InstanceRepository interface {
	Save(instance Instance) error
}

type InstanceRepo struct {
	store store.Store
}

func NewInstanceRepo(s store.Store) *InstanceRepo {
	return &InstanceRepo{store: s}
}

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
