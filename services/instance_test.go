package services

import (
	"testing"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstance(t *testing.T) {
	store := store.New()
	repo := instance.NewInstanceRepo(store)
	service := NewInstanceService(repo)

	ia := instance.New(store, Attributes{PlaybookID: "test", ID: "222"})
	err := service.Create(ia)
	assert.Nil(t, err)
}
