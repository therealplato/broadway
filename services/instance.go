package services

import "github.com/namely/broadway/domain"

// InstanceService definition
type InstanceService struct {
	repo domain.InstanceRepository
}

// NewInstanceService creates a new instance service
func NewInstanceService(r domain.InstanceRepository) *InstanceService {
	return &InstanceService{repo: r}
}

// Create a new instance
func (is *InstanceService) Create(i domain.Instance) error {
	return is.repo.Save(i)
}
