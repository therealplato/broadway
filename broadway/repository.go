package broadway

import "github.com/coreos/etcd/store"

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
