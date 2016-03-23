package broadway

import (
	"testing"

	"github.com/namely/broadway/store"
	"github.com/stretchr/testify/assert"
)

func TestFindByPath(t *testing.T) {
	repo := NewInstanceRepo(store.New())
	i := Instance{PlaybookID: "test", ID: "222"}
	repo.Save(i)

	instance, err := repo.FindByPath(i.Path())
	assert.Nil(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, "test", instance.PlaybookID)
}

func TestFindByPathWhenTheInstanceDoesNotExist(t *testing.T) {
	repo := NewInstanceRepo(store.New())
	i := Instance{PlaybookID: "notcreated", ID: "222"}

	instance, err := repo.FindByPath(i.Path())
	assert.NotNil(t, err)
	assert.Equal(t, "Instance with path: "+i.Path()+" was not found", err.Error())
	assert.Equal(t, "", instance.PlaybookID)
	assert.Equal(t, "", instance.ID)
}
