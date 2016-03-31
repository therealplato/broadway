package services

import (
	"errors"
	"fmt"
	"log"

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

func vars(i *broadway.Instance) map[string]string {
	vs := map[string]string{}
	for k, v := range i.Vars {
		vs[k] = v
	}

	vs["playbook_id"] = i.PlaybookID
	vs["instance_id"] = i.ID
	vs["id"] = i.ID
	vs["instance_status"] = string(i.Status)

	return vs
}

// Deploy deploys a playbook
func (d *DeploymentService) Deploy(instance *broadway.Instance) error {

	playbook, ok := d.playbooks[instance.PlaybookID]
	if !ok {
		return fmt.Errorf("Could not find playbook ID %s while deploying %s\n", instance.PlaybookID, instance.ID)
	}

	deployer := deployment.NewKubernetesDeployment(playbook, vars(instance), d.manifests)

	if instance.Status == broadway.StatusDeploying {
		return errors.New("Instance is being deployed already.")
	}

	if instance.Status == broadway.StatusDeleting {
		return errors.New("Instance is being deleted already.")
	}

	instance.Status = broadway.StatusDeploying
	err := d.repo.Save(instance)
	if err != nil {
		log.Printf("Failed to save instance status Deploying for %s/%s, continuing deployment\n", instance.PlaybookID, instance.ID)
		log.Println(err)
	}

	err = deployer.Deploy()
	if err != nil {
		log.Printf("Deploying %s/%s failed: %s\n", instance.PlaybookID, instance.ID, err.Error())
		instance.Status = broadway.StatusError

		errS := d.repo.Save(instance)
		if errS != nil {
			log.Printf("Failed to save instance status Error for %s/%s\n", instance.PlaybookID, instance.ID)
			log.Println(errS)
			return errS
		}
		return err
	}

	instance.Status = broadway.StatusDeployed
	err = d.repo.Save(instance)
	if err != nil {
		log.Printf("Failed to save instance status Deployed for playbook ID %s, instance %s\n%s\n", instance.PlaybookID, instance.ID, err.Error())
		return err
	}
	return nil

}
