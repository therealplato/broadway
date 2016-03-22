package broadway

import (
	"encoding/json"

	"github.com/namely/broadway/store"
)

// InstanceRepository interface
type InstanceRepository interface {
	Save(instance Instance) error
	FindByPath(path string) (Instance, error)
}

// InstanceRepo handles persistence logic
type InstanceRepo struct {
	store store.Store
}

type NotFound error

// NewInstanceRepo constructor
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

func (ir *InstanceRepo) FindByPath(path string) (Instance, error) {
	instance := Instance{}

	i := ir.store.Value(path)
	err := json.Unmarshal([]byte(i), &instance)
	if err != nil {
		return Instance{}, err
	}
	return instance, nil
}
