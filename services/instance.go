package services

import "github.com/namely/broadway/broadway"

// InstanceService definition
type InstanceService struct {
	repo broadway.InstanceRepository
}

// NewInstanceService creates a new instance service
func NewInstanceService(r broadway.InstanceRepository) *InstanceService {
	return &InstanceService{repo: r}
}

// Create a new instance
func (is *InstanceService) Create(i broadway.Instance) error {
	return is.repo.Save(i)
}
