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

func TestGetDefaultStatus(t *testing.T) {
	s := store.New()
	testcases := []struct {
		attrs Attributes
	}{
		{
			Attributes{PlaybookID: "test", ID: "present"},
		},
	}

	for _, testcase := range testcases {
		i := New(s, &testcase.attrs)
		err := i.Save()
		assert.Nil(t, err)

		status, err := GetStatus(s, testcase.attrs.PlaybookID, testcase.attrs.ID)
		assert.Nil(t, err, "err of GetStatus with good arguments should be nil")
		assert.Equal(t, StatusNew, status, "GetStatus() of newly created instance should be NewStatus")
	}
}

func TestGetStatusSuccess(t *testing.T) {
	s := store.New()

	testcases := []struct {
		attrs Attributes
	}{
		{
			Attributes{PlaybookID: "test", ID: "present", Status: StatusDeployed},
		},
		{
			Attributes{PlaybookID: "test", ID: "present", Status: StatusDeploying},
		},
		{
			Attributes{PlaybookID: "test", ID: "present", Status: StatusError},
		},
		{
			Attributes{PlaybookID: "test", ID: "present", Status: StatusDeleting},
		},
	}

	for _, testcase := range testcases {
		i := New(s, &testcase.attrs)
		err := i.Save()
		assert.Nil(t, err)

		status, err := GetStatus(s, testcase.attrs.PlaybookID, testcase.attrs.ID)
		assert.Nil(t, err, "err of GetStatus with good arguments should be nil")
		assert.Equal(t, testcase.attrs.Status, status, "GetStatus returned unexpected Status")
	}
}
