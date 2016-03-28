package services

import (
	"errors"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/deployment"
	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
	"github.com/namely/broadway/store"
)

// DeploymentService implements the Broadway logic for deployments
type DeploymentService struct {
	repo      broadway.InstanceRepository
	playbooks map[string]*playbook.Playbook
	manifests map[string]*manifest.Manifest
}

// NewDeploymentService creates a new DeploymentService
func NewDeploymentService(s store.Store, ps map[string]*playbook.Playbook, ms map[string]*manifest.Manifest) *DeploymentService {
	return &DeploymentService{
		repo:      broadway.NewInstanceRepo(s),
		playbooks: ps,
		manifests: ms,
	}
}

// Deploy deploys a playbook
func (d *DeploymentService) Deploy(instance *broadway.Instance) error {

	playbook := d.playbooks[instance.PlaybookID]

	deployer := deployment.NewKubernetesDeployment(playbook, instance.Vars, d.manifests)

	if instance.Status == broadway.StatusDeploying {
		return errors.New("Instance is being deployed already.")
	}

	if instance.Status == broadway.StatusDeleting {
		return errors.New("Instance is being deleted already.")
	}

	instance.Status = broadway.StatusDeploying
	d.repo.Save(instance)

	err := deployer.Deploy()
	if err != nil {
		instance.Status = broadway.StatusError
	} else {
		instance.Status = broadway.StatusDeployed
	}
	d.repo.Save(instance)

	return nil
}
