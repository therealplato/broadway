package instance

import (
	"testing"

	"github.com/namely/broadway/broadway"
	"github.com/namely/broadway/services"
	"github.com/namely/broadway/store"

	"github.com/stretchr/testify/assert"
)

func TestSavingInstance(t *testing.T) {
	i := broadway.Instance{PlaybookID: "test", ID: "222"}
	store := store.New()
	service := services.NewInstanceService(store)

	err := service.Create(i)
	assert.Nil(t, err)

	instance, err := service.Show(i.PlaybookID, i.ID)
	assert.Nil(t, err)

	assert.Equal(t, "test", instance.PlaybookID)
}
