package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/store"

	"github.com/stretchr/testify/assert"
)

func TestInstanceCreateWithValidAttributes(t *testing.T) {
	w := httptest.NewRecorder()

	i := map[string]interface{}{
		"playbook_id": "test",
		"id":          "test",
		"vars": map[string]string{
			"version": "ok",
		},
	}

	rbody, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", "/instances", bytes.NewBuffer(rbody))
	req.Header.Add("Content-Type", "application/json")

	mem := store.New()

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code, "Response code should be 201")

	var response instance.InstanceAttributes
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, response.PlaybookId, "test")

	ii, err := instance.Get("test", "test")
	assert.Nil(t, err)
	assert.Equal(t, "test", ii.ID(), "New instance was created")

}

func TestCreateInstanceWithInvalidAttributes(t *testing.T) {
	w := httptest.NewRecorder()

	invalidRequests := map[string]map[string]interface{}{
		"playbook_id": {
			"id": "test",
			"vars": map[string]string{
				"version": "ok",
			},
		},
	}

	for _, i := range invalidRequests {
		rbody, err := json.Marshal(i)
		if err != nil {
			panic(err)
		}

		req, _ := http.NewRequest("POST", "/instances", bytes.NewBuffer(rbody))
		req.Header.Add("Content-Type", "application/json")

		mem := store.New()

		server := New(mem).Handler()
		server.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 422)

		var errorResponse map[string]string

		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		if err != nil {
			panic(err)
		}
		assert.Contains(t, errorResponse["error"], "Unprocessable Entity")
		//assert.Contains(t, errorResponse["error"], field)
	}

}

func TestGetInstanceWithInvalidPath(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/instance/foo/bar", nil)

	mem := store.New()

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)

	var errorResponse map[string]string

	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		panic(err)
	}
	assert.Contains(t, errorResponse["error"], "Not Found")
}

func TestGetInstancesWithUnknownPath(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/instances/foo", nil)

	mem := store.New()

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code, "Response code should be 204 No Content")

	var okResponse []instance.Instance

	err := json.Unmarshal(w.Body.Bytes(), &okResponse)
	if err != nil {
		panic(err)
	}
	if len(okResponse) != 0 {
		t.Errorf("Expected 0 instances matching playbook foo, actual %s\n", len(okResponse))
	}
}
