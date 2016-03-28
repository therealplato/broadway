package services

import (
	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/store"
)

// InstanceService definition
type InstanceService struct {
	repo broadway.InstanceRepository
}

// NewInstanceService creates a new instance service
func NewInstanceService(s store.Store) *InstanceService {
	return &InstanceService{repo: broadway.NewInstanceRepo(s)}
}

// Create a new instance
func (is *InstanceService) Create(i *broadway.Instance) error {
	return is.repo.Save(i)
}

// Show takes playbookID and instanceID and returns the matching Instance, if
// any
func (is *InstanceService) Show(playbookID, ID string) (*broadway.Instance, error) {
	instance, err := is.repo.FindByID(playbookID, ID)
	if err != nil {
		return instance, err
	}
	return instance, nil
}
