package services

import "github.com/namely/broadway/instance"

type InstanceService struct {
	repo instance.InstanceRepository
}

func NewInstanceService(r instance.InstanceRepository) *InstanceService {
	return &InstanceService{repo: r}
}

func (is *InstanceService) Create(i domain.Instance) error {
	return is.repo.Save(i)
}
