package services

import (
	"testing"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstance(t *testing.T) {
	store := store.New()
	service := NewInstanceService(store)

	i := broadway.Instance{PlaybookID: "test", ID: "222"}
	err := service.Create(i)
	assert.Nil(t, err)
	createdInstance, _ := service.Repo.FindByPath(i.Path())
	assert.Equal(t, "test", createdInstance.PlaybookID)
	assert.Equal(t, broadway.StatusNew, createdInstance.Status)
}

func TestGetStatus(t *testing.T) {
	store := store.New()
	service := NewInstanceService(store)

	testcases := []broadway.Instance{
		broadway.Instance{PlaybookID: "test", ID: "present", Status: broadway.StatusDeployed},
		broadway.Instance{PlaybookID: "test", ID: "present", Status: broadway.StatusDeploying},
		broadway.Instance{PlaybookID: "test", ID: "present", Status: broadway.StatusError},
		broadway.Instance{PlaybookID: "test", ID: "present", Status: broadway.StatusDeleting},
	}

	for _, i := range testcases {
		err := service.Create(i)
		assert.Nil(t, err)

		status, err := service.GetStatus(i.Path())
		assert.Nil(t, err, "err of GetStatus with good arguments should be nil")
		assert.Equal(t, i.Status, status, "GetStatus returned unexpected Status")
	}
}

func TestGetStatusFailure(t *testing.T) {
	store := store.New()
	service := NewInstanceService(store)

	i := broadway.Instance{PlaybookID: "notcreated", ID: "222"}

	status, err := service.GetStatus(i.Path())
	assert.NotNil(t, err)
	assert.Equal(t, broadway.StatusNew, status, "status of GetStatus(missing) should be StatusNew")
}
