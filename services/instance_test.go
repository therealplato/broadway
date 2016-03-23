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

func TestGetStatusFailure(t *testing.T) {
	store := store.New()
	service := NewInstanceService(store)

	i := broadway.Instance{PlaybookID: "notcreated", ID: "222"}

	status, err := service.GetStatus(i.Path())
	assert.NotNil(t, err)
	assert.Equal(t, broadway.StatusNew, status, "status of GetStatus(missing) should be StatusNew")
}

func TestGetStatus(t *testing.T) {
}
