package services

import (
	"testing"

	"github.com/namely/broadway/domain"
	"github.com/namely/broadway/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstance(t *testing.T) {
	store := store.New()
	repo := domain.NewInstanceRepo(store)
	service := NewInstanceService(repo)

	i := domain.Instance{PlaybookID: "test", ID: "222"}
	err := service.Create(i)
	assert.Nil(t, err)
}
