package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/store"

	"github.com/stretchr/testify/assert"
)

func TestServerNew(t *testing.T) {
	_, exists := os.LookupEnv(slackTokenENV)
	if exists {
		t.Fatalf("Found existing $%s. Skipping tests to avoid changing it...", slackTokenENV)
	}

	testToken := "BroadwayTestToken"

	err := os.Setenv(slackTokenENV, testToken)
	if err != nil {
		t.Fatal(err)
	}
	actualToken, exists := os.LookupEnv(slackTokenENV)
	assert.True(t, exists, "Expected ENV to exist")
	assert.Equal(t, testToken, actualToken, "Unexpected ENV value")

	mem := store.New()

	s := New(mem)
	assert.Equal(t, testToken, s.slackToken, "Expected server.slackToken to match existing ENV value")

	err = os.Unsetenv(slackTokenENV)
	if err != nil {
		t.Fatal(err)
	}
	actualToken, exists = os.LookupEnv(slackTokenENV)
	assert.False(t, exists, "Expected ENV to not exist")
	assert.Equal(t, "", actualToken, "Unexpected ENV value")
	s = New(mem)
	assert.Equal(t, "", s.slackToken, "Expected server.slackToken to be empty string for missing ENV value")

}

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
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", "/instances", bytes.NewBuffer(rbody))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	mem := store.New()

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code, "Response code should be 201")

	var response instance.Attributes
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, response.PlaybookID, "test")

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
			t.Error(err)
			return
		}

		req, err := http.NewRequest("POST", "/instances", bytes.NewBuffer(rbody))
		if err != nil {
			t.Error(err)
			return
		}

		req.Header.Add("Content-Type", "application/json")

		mem := store.New()

		server := New(mem).Handler()
		server.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 422)

		var errorResponse map[string]string

		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Contains(t, errorResponse["error"], "Unprocessable Entity")
		//assert.Contains(t, errorResponse["error"], field)
	}

}

func TestGetInstanceWithValidPath(t *testing.T) {
	w := httptest.NewRecorder()
	mem := store.New()

	i := instance.New(mem, &instance.Attributes{
		PlaybookID: "foo",
		ID:         "doesExist",
	})
	err := i.Save()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("GET", "/instance/foo/doesExist", nil)
	if err != nil {
		t.Error(err)
		return
	}

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	var iResponse map[string]string

	err = json.Unmarshal(w.Body.Bytes(), &iResponse)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Contains(t, iResponse["id"], "doesExist")
}

func TestGetInstanceWithInvalidPath(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/instance/foo/bar", nil)
	if err != nil {
		t.Error(err)
		return
	}

	mem := store.New()

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)

	var errorResponse map[string]string

	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Contains(t, errorResponse["error"], "Not Found")
}

func TestGetInstancesWithFullPlaybook(t *testing.T) {
	w := httptest.NewRecorder()
	mem := store.New()

	testInstance1 := instance.New(mem, &instance.Attributes{
		PlaybookID: "testPlaybookFull",
		ID:         "testInstance1",
	})
	err := testInstance1.Save()
	if err != nil {
		t.Error(err)
		return
	}
	testInstance2 := instance.New(mem, &instance.Attributes{
		PlaybookID: "testPlaybookFull",
		ID:         "testInstance2",
	})
	err = testInstance2.Save()
	if err != nil {
		t.Error(err)
		return
	}
	req, err := http.NewRequest("GET", "/instances/testPlaybookFull", nil)
	if err != nil {
		t.Error(err)
		return
	}

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200 OK")

	log.Println(w.Body.String())
	var okResponse []instance.Attributes

	err = json.Unmarshal(w.Body.Bytes(), &okResponse)
	if err != nil {
		t.Error(err)
		return
	}
	if len(okResponse) != 2 {
		t.Errorf("Expected 2 instances matching playbook testPlaybookFull, actual %v\n", len(okResponse))
	}
}

func TestGetInstancesWithEmptyPlaybook(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/instances/testPlaybookEmpty", nil)
	if err != nil {
		t.Error(err)
		return
	}

	mem := store.New()

	server := New(mem).Handler()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code, "Response code should be 204 No Content")

	var okResponse []instance.Attributes

	err = json.Unmarshal(w.Body.Bytes(), &okResponse)
	if err != nil {
		t.Error(err)
		return
	}
	if len(okResponse) != 0 {
		t.Errorf("Expected 0 instances matching playbook testPlaybookEmpty, actual %v\n", len(okResponse))
	}
}
