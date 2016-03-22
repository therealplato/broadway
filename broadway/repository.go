package broadway

import "github.com/namely/broadway/store"

// InstanceRepository interface
type InstanceRepository interface {
	Save(instance Instance) error
}

// InstanceRepo handles persistence logic
type InstanceRepo struct {
	store store.Store
}

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
