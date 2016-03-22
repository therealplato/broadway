package domain

type InstanceRepository interface {
	Save(instance Instance) error
}

type InstanceRepo struct {
	store store.Store
}

func NewInstanceRepo(s store.Store) *InstanceRepo {
	return &InstanceRepo{store: s}
}

func (ir *InstanceRepo) Save(instance ) error {
	encoded, err := i.Attributes().JSON()
	if err != nil {
		return err
	}
	err = ir.store.SetValue(i., encoded)
	if err != nil {
		return err
	}
	return nil
}
