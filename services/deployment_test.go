package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/namely/broadway/store/etcdstore"
	"github.com/namely/broadway/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/namely/broadway/deployment"
	"github.com/namely/broadway/instance"
)

func init() {
	fmt.Println("82481234891238429318491238492138498231984912384912384")
	fmt.Printf("%+v\n", ServicesTestCfg)
	etcdstore.Setup(testutils.TestCfg)
	deployment.Setup(testutils.TestCfg)
}

func TestDeployment(t *testing.T) {
	nt := newNotificationTestHelper()
	defer nt.Close()
	manifests, err := NewManifestService(ServicesTestCfg).LoadManifestFolder()
	if err != nil {
		panic(err)
	}

	playbooks, err := deployment.LoadPlaybookFolder(ServicesTestCfg.PlaybooksPath)
	fmt.Printf("%+v\n", playbooks)
	if err != nil {
		panic(err)
	}

	ds := NewDeploymentService(ServicesTestCfg, etcdstore.New(), playbooks, manifests)

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
		err = ds.DeployAndNotify(c.Instance)
		assert.Equal(t, c.Error, err)
		assert.EqualValues(t, c.Expected, c.Instance.Status)
	}
}

func TestCustomDeploymentNotification(t *testing.T) {
	nt := newNotificationTestHelper()
	defer nt.Close()
	manifests, err := NewManifestService(ServicesTestCfg).LoadManifestFolder()
	if err != nil {
		panic(err)
	}

	playbooks, err := deployment.LoadPlaybookFolder(ServicesTestCfg.PlaybooksPath)
	if err != nil {
		panic(err)
	}
	service := NewDeploymentService(ServicesTestCfg, etcdstore.New(), playbooks, manifests)

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
