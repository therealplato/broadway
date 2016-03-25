package services

import (
	"log"

	"github.com/namely/broadway/deployment"
	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
)

// DeploymentService mediates between high-level layers that want to deploy
// things and the low-level implementation of that deployment
type DeploymentService struct {
	deployer deployment.Deployer
}

// Deploy deploys a playbook
func (d *DeploymentService) Deploy(p playbook.Playbook, vars map[string]string) error {
	return d.deployer.Deploy(p, vars)
}

// DefaultDeployer implements the Deployer interface
type DefaultDeployer struct{}

// Deploy finds the playbook's manifests and executes them with Kubernetes.
func (dd *DefaultDeployer) Deploy(p playbook.Playbook, vars map[string]string) error {
	MS := NewManifestService()
	for _, t := range p.Tasks {
		pod, manifests, err := MS.LoadTask(t)
		_ = pod
		_ = err
		// execute templates with vars

		// give templated manifest to kubernetes
		go func(ms []manifest.Manifest) {
			log.Printf("Deploying: %+v\n", ms)
		}(manifests)
	}
	return nil
}
