package services

import (
	"testing"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstance(t *testing.T) {
	store := store.New()
	repo := broadway.NewInstanceRepo(store)
	service := NewInstanceService(repo)

	i := broadway.Instance{PlaybookID: "test", ID: "222"}
	err := service.Create(i)
	assert.Nil(t, err)
	createdInstance, _ := repo.FindByPath(i.Path())
	assert.Equal(t, "test", createdInstance.PlaybookID)
}
