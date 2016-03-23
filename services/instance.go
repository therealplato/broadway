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
