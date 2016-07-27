package services

import (
	"errors"
	"testing"

	"github.com/namely/broadway/store/etcdstore"
	"github.com/stretchr/testify/assert"

	"github.com/namely/broadway/deployment"
	"github.com/namely/broadway/env"
	"github.com/namely/broadway/instance"
)

func TestDeployment(t *testing.T) {
	nt := newNotificationTestHelper()
	defer nt.Close()
	manifests, err := NewManifestService(env.ManifestsPath).LoadManifestFolder()
	if err != nil {
		panic(err)
	}

	playbooks, err := deployment.LoadPlaybookFolder(env.PlaybooksPath)
	if err != nil {
		panic(err)
	}

	service := NewDeploymentService(etcdstore.New(), playbooks, manifests)

	cases := []struct {
		Name     string
		Instance *instance.Instance
		Error    error
		Expected instance.Status
	}{
		{
			Name: "Good Playbook",
			Instance: &instance.Instance{
				PlaybookID: "hello",
				ID:         "test",
				Vars: map[string]string{
					"version": "test",
				},
			},
			Error:    nil,
			Expected: instance.StatusDeployed,
		}, {
			Name: "Playbook con Mal Pod",
			Instance: &instance.Instance{
				PlaybookID: "goodbye",
				ID:         "another",
				Vars: map[string]string{
					"version": "another",
				},
			},
			Error:    errors.New("Setup pod failed!"),
			Expected: instance.StatusError,
		},
	}

	for _, c := range cases {
		err = service.DeployAndNotify(c.Instance)
		assert.Equal(t, c.Error, err)
		assert.EqualValues(t, c.Expected, c.Instance.Status)
	}
}

func TestCustomDeploymentNotification(t *testing.T) {
	nt := newNotificationTestHelper()
	defer nt.Close()
	manifests, err := NewManifestService(env.ManifestsPath).LoadManifestFolder()
	if err != nil {
		panic(err)
	}

	playbooks, err := deployment.LoadPlaybookFolder(env.PlaybooksPath)
	if err != nil {
		panic(err)
	}
	service := NewDeploymentService(etcdstore.New(), playbooks, manifests)

	i := &instance.Instance{
		PlaybookID: "messagesplaybook",
		ID:         "test",
		Vars: map[string]string{
			"version": "test",
		},
	}

	err = service.DeployAndNotify(i)
	assert.Equal(t, nil, err)
	assert.Contains(t, nt.requestBody, "custom deployed")
	assert.Contains(t, nt.requestBody, "messagesplaybook/test")
}
