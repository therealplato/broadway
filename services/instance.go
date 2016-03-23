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
	r := broadway.NewInstanceRepo(s)
	return &InstanceService{repo: r}
}

// Create a new instance
func (is *InstanceService) Create(i broadway.Instance) error {
	return is.repo.Save(i)
}

func (is *InstanceService) Show(playbookId, id string) (broadway.Instance, error) {
	instance, err := is.repo.FindById(playbookId, id)
	if err != nil {
	}
	return instance, nil
}
