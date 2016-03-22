package services

import "github.com/namely/broadway/domain"

type InstanceService struct {
	repo domain.InstanceRepository
}

func NewInstanceService(r domain.InstanceRepository) *InstanceService {
	return &InstanceService{repo: r}
}

func (is *InstanceService) Create(i domain.Instance) error {
	return is.repo.Save(i)
}
