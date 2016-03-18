package instance

import (
	"testing"

	"github.com/namely/broadway/store"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusFailure(t *testing.T) {
	s := store.New()

	status, err := GetStatus(s, "test", "missing")
	assert.Equal(t, StatusNew, status, "status of GetStatus(missing) should be StatusNew")
	assert.Equal(t, "Instance does not exist.", err.Error(), "wrong err: GetStatus with bad arguments")
}

func TestGetStatusSuccess(t *testing.T) {
	s := store.New()
	i := New(s, &Attributes{PlaybookID: "test", ID: "present"})
	err := i.Save()
	assert.Nil(t, err)

	status, err := GetStatus(s, "test", "present")
	assert.Nil(t, err, "err of GetStatus with good arguments should be nil")
	assert.Equal(t, StatusNew, status, "Status of GetStatus(new) should be StatusNew")
}
