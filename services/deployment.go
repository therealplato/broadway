package services

import (
	"github.com/namely/broadway/deployment"
	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
)

// DeploymentService implements the Broadway logic for deployments
type DeploymentService struct {
	deployer  deployment.Deployer
	playbooks map[string]*playbook.Playbook
	manifests map[string]*manifests.Manifest
}

// NewDeploymentService creates a new DeploymentService
func NewDeploymentService(deployer deployment.Deployer, playbooks map[string]*playbook.Playbook, manifests map[string]*manifest.Manifest) *DeploymentService {
	return &DeploymentService{
		deployer:  deployer,
		playbooks: playbooks,
		manifests: manifests,
	}
}

// Deploy deploys a playbook
func (d *DeploymentService) Deploy(instance) error {

	//if instance.Status != Deploying .. {
	//	instance.Status = Deploying
	//	instance.Save

	//	err := d.deployer.Deploy(p, instance.Vars)
	//	if err != nil {
	//		instance.Status = error
	//		instance.save
	//	} else {
	//		instance.Status = deployed
	//		instance.save
	//	}

	//}
}
