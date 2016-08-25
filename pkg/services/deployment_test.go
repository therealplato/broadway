package services

import (
	"testing"
	"time"

	"github.com/namely/broadway/pkg/store/etcdstore"
	"github.com/namely/broadway/pkg/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/namely/broadway/pkg/deployment"
	"github.com/namely/broadway/pkg/instance"
)

func init() {
	etcdstore.Setup(testutils.TestCfg)
	deployment.Setup(testutils.TestCfg)
}

func TestRemoveExpiredInstancest(t *testing.T) {
	nt := newNotificationTestHelper()
	defer nt.Close()
	manifests, err := deployment.LoadManifestFolder(ServicesTestCfg.ManifestsPath, ServicesTestCfg.ManifestsExtension)
	if err != nil {
		panic(err)
	}

	playbooks, err := deployment.LoadPlaybookFolder(ServicesTestCfg.PlaybooksPath)
	if err != nil {
		panic(err)
	}

	ds := NewDeploymentService(ServicesTestCfg, etcdstore.New(), playbooks, manifests)

	cases := []struct {
		Scenario       string
		Instance       *instance.Instance
		CurrentDate    time.Time
		ExpirationDate time.Time
		Error          error
	}{
		{
			Scenario:       "RemoveExpiredInstances: When an instance just expired",
			CurrentDate:    time.Date(2016, 8, 05, 0, 00, 00, 651387237, time.UTC),
			ExpirationDate: time.Date(2016, 8, 10, 0, 00, 00, 651387237, time.UTC),
			Instance: &instance.Instance{
				PlaybookID: "hello",
				ID:         "anothertest",
				Path: instance.Path{
					RootPath:   ServicesTestCfg.EtcdPath,
					PlaybookID: "hello",
					ID:         "anothertest",
				},
			},
			Error: nil,
		},
	}

	s := etcdstore.New()
	for _, c := range cases {
		c.Instance.ExpiredAt = instance.NewExpiredAt(ServicesTestCfg.InstanceExpirationDays, c.CurrentDate).Unix()
		err := instance.Save(s, c.Instance)
		assert.Nil(t, err, c.Scenario)

		err = ds.DeployAndNotify(c.Instance)
		assert.Nil(t, err, c.Scenario)

		err = ds.RemoveExpiredInstances(c.ExpirationDate)

		ii, err := instance.FindByPath(s, c.Instance.Path)
		assert.Equal(t, c.Error, err, c.Scenario)
		assert.Equal(t, instance.StatusDeleting, ii.Status, c.Scenario)
	}
}

func TestDeployment(t *testing.T) {
	nt := newNotificationTestHelper()
	defer nt.Close()
	manifests, err := deployment.LoadManifestFolder(ServicesTestCfg.ManifestsPath, ServicesTestCfg.ManifestsExtension)
	if err != nil {
		panic(err)
	}

	playbooks, err := deployment.LoadPlaybookFolder(ServicesTestCfg.PlaybooksPath)
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
	manifests, err := deployment.LoadManifestFolder(ServicesTestCfg.ManifestsPath, ServicesTestCfg.ManifestsExtension)
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
