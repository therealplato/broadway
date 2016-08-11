package instance

import (
	"errors"
	"testing"
	"time"

	"github.com/namely/broadway/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestNewExpiredAt(t *testing.T) {
	testcases := []struct {
		Scenario     string
		DaysToExpire int
		CurrentTime  time.Time
		ExpectedTime time.Time
	}{
		{
			Scenario:     "NewExpiredAt: Build a new expired at type",
			DaysToExpire: 5,
			CurrentTime:  time.Date(2016, 8, 5, 00, 00, 00, 651387237, time.UTC),
			ExpectedTime: time.Date(2016, 8, 10, 00, 00, 00, 651387237, time.UTC),
		},
		{
			Scenario:     "NewExpiredAt: When day is 29",
			DaysToExpire: 5,
			CurrentTime:  time.Date(2016, 8, 29, 00, 00, 00, 651387237, time.UTC),
			ExpectedTime: time.Date(2016, 9, 3, 00, 00, 00, 651387237, time.UTC),
		},
	}

	for _, tc := range testcases {
		expiredAt := NewExpiredAt(tc.DaysToExpire, tc.CurrentTime)
		assert.Equal(t, tc.ExpectedTime.Year(), expiredAt.Year(), tc.Scenario)
		assert.Equal(t, tc.ExpectedTime.Month(), expiredAt.Month(), tc.Scenario)
		assert.Equal(t, tc.ExpectedTime.Day(), expiredAt.Day(), tc.Scenario)
	}
}

func TestFindByPath(t *testing.T) {
	testcases := []struct {
		Scenario           string
		Path               Path
		Store              store.Store
		ExpectedPlaybookID string
		ExpectedError      error
	}{
		{
			Scenario: "When the instance is properly save",
			Path:     Path{"etcdPath", "test", "id"},
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":"test", "id": "id", "status": "deployed"}`
				},
			},
			ExpectedPlaybookID: "test",
			ExpectedError:      nil,
		},
		{
			Scenario: "When the instance was not properly save",
			Path:     Path{"etcdPath", "test", "id"},
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":}`
				},
			},
			ExpectedPlaybookID: "",
			ExpectedError:      ErrMalformedSaveData,
		},
		{
			Scenario: "When the instance does not exist",
			Path:     Path{"etcdPath", "test", "id"},
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return ""
				},
			},
			ExpectedPlaybookID: "",
			ExpectedError:      NotFoundError("etcdPath/instances/test/id"),
		},
	}

	for _, tc := range testcases {
		returnedInstance, err := FindByPath(tc.Store, tc.Path)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
		if err == nil {
			assert.Equal(t, tc.ExpectedPlaybookID, returnedInstance.PlaybookID)
		}
	}
}

func TestFindByPlaybookID(t *testing.T) {
	testcases := []struct {
		Scenario          string
		Store             store.Store
		PlaybookPath      PlaybookPath
		ExpectedInstances []*Instance
		ExpectedError     error
	}{
		{
			Scenario: "When instances exist in the store",
			Store: &store.FakeStore{
				MockValues: func(string) map[string]string {
					return map[string]string{
						"rootPath/instances/test":  `{"playbook_id": "test", "id": "id", "status": "deployed"}`,
						"rootPath/instances/test1": `{"playbook_id": "test1", "id": "id", "status": "deployed"}`,
					}
				},
			},
			PlaybookPath: PlaybookPath{"rootPath", "test"},
			ExpectedInstances: []*Instance{
				&Instance{PlaybookID: "test", ID: "id", Status: StatusDeployed},
				&Instance{PlaybookID: "test1", ID: "id", Status: StatusDeployed},
			},
			ExpectedError: nil,
		},
		{
			Scenario: "When instances does not exist in the store",
			Store: &store.FakeStore{
				MockValues: func(string) map[string]string {
					return nil
				},
			},
			PlaybookPath:      PlaybookPath{"rootPath", "test"},
			ExpectedInstances: []*Instance{},
			ExpectedError:     nil,
		},
		{
			Scenario: "When the data is malformed",
			Store: &store.FakeStore{
				MockValues: func(string) map[string]string {
					return map[string]string{
						"rootPath/instances/test":  `{"playbook_id": "test", "id": "id", "status": "deployed"}`,
						"rootPath/instances/test1": `{"playbook_id":`,
					}
				},
			},
			PlaybookPath:      PlaybookPath{"rootPath", "test"},
			ExpectedInstances: nil,
			ExpectedError:     ErrMalformedSaveData,
		},
	}

	for _, tc := range testcases {
		instances, err := FindByPlaybookPath(tc.Store, tc.PlaybookPath)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
		if err == nil {
			assert.Equal(t, tc.ExpectedInstances, instances, tc.Scenario)
		}
	}
}

func TestSave(t *testing.T) {
	testcases := []struct {
		Scenario      string
		Store         store.Store
		Instance      *Instance
		ExpectedError error
	}{
		{
			Scenario: "Save: When successfully save in store",
			Store: &store.FakeStore{
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Instance:      &Instance{PlaybookID: "playbookID", ID: "id"},
			ExpectedError: nil,
		},
		{
			Scenario: "Save: When successfully with ExpiredAt set save in store",
			Store: &store.FakeStore{
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Instance:      &Instance{PlaybookID: "playbookID", ID: "id", ExpiredAt: time.Now().Unix()},
			ExpectedError: nil,
		},
	}
	for _, tc := range testcases {
		err := Save(tc.Store, tc.Instance)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
	}
}

func TestDelete(t *testing.T) {
	testcases := []struct {
		Scenario      string
		Store         store.Store
		Path          Path
		ExpectedError error
	}{
		{
			Scenario: "When successfully deleted from store",
			Store: &store.FakeStore{
				MockDelete: func(string) error {
					return nil
				},
			},
			Path:          Path{"rootPath", "playbookId", "id"},
			ExpectedError: nil,
		},
	}

	for _, tc := range testcases {
		err := Delete(tc.Store, tc.Path)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
	}
}

func TestLock(t *testing.T) {
	testcases := []struct {
		Scenario         string
		Store            store.Store
		Path             Path
		ExpectedError    error
		ExpectedInstance *Instance
	}{
		{
			Scenario: "Lock: when an instance exist",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":"test", "id": "id", "status": "deployed", "lock": false}`
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    nil,
			ExpectedInstance: &Instance{PlaybookID: "test", ID: "id", Status: StatusDeployed, Lock: true},
		},
		{
			Scenario: "Lock: when an instance does not exist",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return ``
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    NotFoundError("rootPath/instances/test/id"),
			ExpectedInstance: nil,
		},
		{
			Scenario: "Lock: when instance have malformed data saved",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{play}`
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    ErrMalformedSaveData,
			ExpectedInstance: nil,
		},
		{
			Scenario: "Lock: when saving the instance with the new status failed",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":"test", "id": "id", "status": "deployed"}`
				},
				MockSetValue: func(string, string) error {
					return errors.New("Failed to save the instance")
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    errors.New("Failed to save the instance"),
			ExpectedInstance: nil,
		},
	}

	for _, tc := range testcases {
		instance, err := Lock(tc.Store, tc.Path)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
		assert.Equal(t, tc.ExpectedInstance, instance, tc.Scenario)
	}
}

func TestUnlock(t *testing.T) {
	testcases := []struct {
		Scenario         string
		Store            store.Store
		Path             Path
		ExpectedError    error
		ExpectedInstance *Instance
	}{
		{
			Scenario: "Unlock: when an instance exist and is locked",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":"test", "id": "id", "status": "deployed", "lock": true}`
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    nil,
			ExpectedInstance: &Instance{PlaybookID: "test", ID: "id", Status: StatusDeployed, Lock: false},
		},
		{
			Scenario: "Unlock: when an instance does not exist",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return ``
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    NotFoundError("rootPath/instances/test/id"),
			ExpectedInstance: nil,
		},
		{
			Scenario: "Unlock: when instance have malformed data saved",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{play}`
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    ErrMalformedSaveData,
			ExpectedInstance: nil,
		},
		{
			Scenario: "Unlock: when saving the instance failed",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":"test", "id": "id", "status": "deployed", "lock": true}`
				},
				MockSetValue: func(string, string) error {
					return errors.New("Failed to save the instance")
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    errors.New("Failed to save the instance"),
			ExpectedInstance: nil,
		},
		{
			Scenario: "Unlock: when instance is not locked",
			Store: &store.FakeStore{
				MockValue: func(path string) string {
					return `{"playbook_id":"test", "id": "id", "status": "deployed", "lock": false}`
				},
				MockSetValue: func(string, string) error {
					return nil
				},
			},
			Path:             Path{"rootPath", "test", "id"},
			ExpectedError:    NotLockedStatusError("rootPath/instances/test/id"),
			ExpectedInstance: nil,
		},
	}

	for _, tc := range testcases {
		instance, err := Unlock(tc.Store, tc.Path)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
		assert.Equal(t, tc.ExpectedInstance, instance, tc.Scenario)
	}
}

func TestAllDeployedAndExpired(t *testing.T) {
	testcases := []struct {
		Scenario          string
		Path              string
		Store             store.Store
		ExpirationDate    time.Time
		ExpectedInstances []*Instance
		ExpectedError     error
	}{
		{
			Scenario:       "AllDeployedAndExpired: When instance was deployed and expired",
			Path:           "broadwaytest/instances",
			ExpirationDate: time.Date(2016, 8, 5, 00, 00, 00, 651387237, time.UTC),
			Store: &store.FakeStore{
				MockValues: func(path string) map[string]string {
					return map[string]string{
						"etcdPath/instances": `{"playbook_id":"test", "id": "id", "status": "deployed", "expired_at": 10}`,
					}
				},
			},
			ExpectedInstances: []*Instance{
				&Instance{PlaybookID: "test", ID: "id", Status: StatusDeployed, ExpiredAt: 10},
			},
			ExpectedError: nil,
		},
		{
			Scenario:       "AllDeployedAndExpired: When instance was deployed and expired today",
			Path:           "broadwaytest/instances",
			ExpirationDate: time.Date(2016, 8, 5, 00, 00, 00, 651387237, time.UTC),
			Store: &store.FakeStore{
				MockValues: func(path string) map[string]string {
					return map[string]string{
						"etcdPath/instances": `{"playbook_id":"test", "id": "id", "status": "deployed", "expired_at": 1470355200}`,
					}
				},
			},
			ExpectedInstances: []*Instance{
				&Instance{PlaybookID: "test", ID: "id", Status: StatusDeployed, ExpiredAt: 1470355200},
			},
			ExpectedError: nil,
		},
	}

	for _, tc := range testcases {
		instances, err := AllDeployedAndExpired(tc.Store, tc.Path, tc.ExpirationDate)
		assert.Equal(t, tc.ExpectedError, err, tc.Scenario)
		assert.Equal(t, tc.ExpectedInstances, instances, tc.Scenario)
	}
}
