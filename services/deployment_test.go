package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/playbook"
	"github.com/namely/broadway/store"
)

func TestDeployment(t *testing.T) {
	manifests, err := NewManifestService("../examples/manifests/").LoadManifestFolder()
	if err != nil {
		panic(err)
	}

	playbooks, err := playbook.LoadPlaybookFolder("../examples/playbooks")
	if err != nil {
		panic(err)
	}

	service := NewDeploymentService(store.New(), playbooks, manifests)

	i := &broadway.Instance{
		PlaybookID: "hello",
		ID:         "test",
		Vars: map[string]string{
			"version": "test",
		},
	}

	err = service.Deploy(i)
	assert.Nil(t, err)
	assert.EqualValues(t, broadway.StatusDeployed, i.Status)
}
