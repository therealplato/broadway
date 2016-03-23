package services

import (
	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/store"
)

// InstanceService definition
type InstanceService struct {
	Repo broadway.InstanceRepository
}

// NewInstanceService creates a new instance service
func NewInstanceService(s store.Store) *InstanceService {
	r := broadway.NewInstanceRepo(s)
	return &InstanceService{Repo: r}
}

// Create a new instance
func (is *InstanceService) Create(i broadway.Instance) error {
	return is.Repo.Save(i)
}

func (is *InstanceService) GetStatus(path string) (broadway.Status, error) {
	instance, err := is.Repo.FindByPath(path)
	if err != nil {
		return broadway.StatusNew, err
	}
	return instance.Status, nil
}
