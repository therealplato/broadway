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
