package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/coreos/etcd/store"
	"github.com/namely/broadway/cfg"
	"github.com/namely/broadway/instance"
	"github.com/namely/broadway/services"
	"github.com/namely/broadway/store/etcdstore"
	"github.com/namely/broadway/testutils"

	"github.com/stretchr/testify/assert"
)

var testToken = "BroadwayTestToken"
var testCommonCfg = cfg.CommonCfgType{}
var testServerCfg = cfg.ServerCfgType{
	AuthBearerToken: "testtoken",
	SlackToken:      testToken,
	ManifestsPath:   "examples/manifests",
	PlaybooksPath:   "examples/playbooks",
}

func makeRequest(s *Server, req *http.Request, w *httptest.ResponseRecorder) {
	s.Init()
	s.Handler().ServeHTTP(w, req)
}

func auth(cfg cfg.ServerCfgType, req *http.Request) *http.Request {
	req.Header.Set("Authorization", "Bearer "+cfg.AuthBearerToken)
	return req
}

func TestServerNew(t *testing.T) {
	_, s, _ := helperSetupServer(testServerCfg)
	assert.Equal(t, testServerCfg, s.Cfg, "Expected server.Cfg to match passed in config")
}

func TestAuthFailure(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer faketoken")
	w, _, e := helperSetupServer(testServerCfg)
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected POST / with wrong auth token to be 401")
	assert.Contains(t, w.Body.String(), "Authorization")
}

func TestAuthSuccess(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer testtoken")
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected POST / with correct auth token to be 200")
}

func TestInstanceCreateWithValidAttributes(t *testing.T) {
	i := map[string]interface{}{
		"playbook_id": "helloplaybook",
		"id":          "TestInstanceCreateWithValidAttributes",
		"vars": map[string]string{
			"word": "gorilla",
		},
	}
	rbody := testutils.JSONFromMap(t, i)
	req, w := testutils.PostRequest(t, "/instances", rbody)
	req = auth(testServerCfg, req)
	server := New(store.New(), testCommonCfg, testServerCfg)
	makeRequest(server, req, w)
	assert.Equal(t, http.StatusCreated, w.Code, "Response code should be 201")
}

func TestCreateInstanceWithInvalidAttributes(t *testing.T) {
	invalidRequests := map[string]map[string]interface{}{
		"playbook_id": {
			"id": "test",
			"vars": map[string]string{
				"version": "ok",
			},
		},
	}

	for _, i := range invalidRequests {
		rbody := testutils.JSONFromMap(t, i)
		req, w := testutils.PostRequest(t, "/instances", rbody)
		req = auth(testServerCfg, req)
		server := New(store.New(), testCommonCfg, testServerCfg)
		makeRequest(server, req, w)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected POST /instances with wrong attributes to be 400")
	}
}

func TestGetInstanceWithValidPath(t *testing.T) {
	st := store.New()
	i := &instance.Instance{PlaybookID: "helloplaybook", ID: "TestGetInstanceWithValidPath"}
	service := services.NewInstanceService(st)
	_, err := service.CreateOrUpdate(i)
	if err != nil {
		t.Log(err.Error())
	}

	req, w := testutils.GetRequest(t, "/instance/helloplaybook/TestGetInstanceWithValidPath")
	req = auth(testServerCfg, req)
	server := New(st, testCommonCfg, testServerCfg)
	makeRequest(server, req, w)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetInstanceWithInvalidPath(t *testing.T) {
	req, w := testutils.GetRequest(t, "/instance/vanished/TestGetInstanceWithInvalidPath")
	req = auth(testServerCfg, req)
	server := New(store.New(), testCommonCfg, testServerCfg)
	makeRequest(server, req, w)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetInstancesWithFullPlaybook(t *testing.T) {
	testInstance1 := &instance.Instance{PlaybookID: "helloplaybook", ID: "TestGetInstancesWithFullPlaybook1"}
	testInstance2 := &instance.Instance{PlaybookID: "helloplaybook", ID: "TestGetInstancesWithFullPlaybook2"}
	service := services.NewInstanceService(etcdstore.New())
	_, err := service.CreateOrUpdate(testInstance1)
	_, err = service.CreateOrUpdate(testInstance2)
	if err != nil {
		t.Log(err.Error())
	}

	req, w := testutils.GetRequest(t, "/instances/helloplaybook")
	req = auth(testServerCfg, req)
	server := New(store.New(), testCommonCfg, testServerCfg)
	makeRequest(server, req, w)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200 OK")
}

func TestGetStatusFailures(t *testing.T) {
	invalidRequests := []struct {
		method  string
		path    string
		errCode int
		errMsg  string
	}{
		{
			"GET",
			"/status/helloplaybook/TestGetStatusFailures",
			404,
			"Not Found",
		},
	}

	for _, i := range invalidRequests {
		req, w := testutils.GetRequest(t, i.path)
		req = auth(testServerCfg, req)
		server := New(store.New(), testCommonCfg, testServerCfg)
		makeRequest(server, req, w)

		assert.Equal(t, i.errCode, w.Code)

		var errorResponse map[string]string

		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Nil(t, err)
		assert.Contains(t, errorResponse["error"], i.errMsg)
	}

}
func TestGetStatusWithGoodPath(t *testing.T) {
	testInstance1 := &instance.Instance{
		PlaybookID: "helloplaybook",
		ID:         "TestGetStatusWithGoodPath",
		Status:     instance.StatusDeployed}
	is := services.NewInstanceService(etcdstore.New())
	_, err := is.CreateOrUpdate(testInstance1)
	if err != nil {
		t.Fatal(err)
	}
	req, w := testutils.GetRequest(t, "/status/helloplaybook/TestGetStatusWithGoodPath")
	req = auth(testServerCfg, req)
	server := New(store.New(), testCommonCfg, testServerCfg)
	makeRequest(server, req, w)

	assert.Equal(t, http.StatusOK, w.Code)

	var statusResponse map[string]string

	err = json.Unmarshal(w.Body.Bytes(), &statusResponse)
	assert.Nil(t, err)
	assert.Contains(t, statusResponse["status"], "deployed")
}

func helperSetupServer(cfg cfg.ServerCfgType) (*httptest.ResponseRecorder, *Server, http.Handler) {
	w := httptest.NewRecorder()
	mem := store.New()
	s := New(mem, testCommonCfg, cfg)
	return w, s, s.Handler()
}

func TestGetCommand400(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, err := http.NewRequest("GET", "/command", nil)
	if err != nil {
		t.Fatal(err)
	}
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected GET /command to be 400")
}

func TestGetCommand200(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, _ := http.NewRequest("GET", "/command?ssl_check=1", nil)
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected GET /command?ssl_check=1 to be 200")
}
func TestPostCommandMissingToken(t *testing.T) {
	formBytes := bytes.NewBufferString("not a form")
	req, _ := http.NewRequest("POST", "/command", formBytes)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testCfg := cfg.ServerCfgType{SlackToken: testToken}
	w, _, e := helperSetupServer(testCfg)
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected POST /command with bad body to be 401")
}
func TestPostCommandWrongToken(t *testing.T) {
	testCfg := cfg.ServerCfgType{SlackToken: testToken}
	w, _, e := helperSetupServer(testCfg)
	req, _ := http.NewRequest("POST", "/command", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	form := url.Values{}
	form.Set("token", "wrongtoken")
	req.PostForm = form
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected POST /command with wrong token to be 401")
}
func TestPostCommandHelp(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, _ := http.NewRequest("POST", "/command", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	form := url.Values{}
	form.Set("token", testToken)
	form.Set("command", "/broadway")
	form.Set("text", "help")
	req.PostForm = form

	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected /broadway help to be 200")
	assert.Contains(t, w.Body.String(), "deploy", "Expected help message to contain deploy")
	assert.Contains(t, w.Body.String(), "info", "Expected help message to contain info")
	assert.Contains(t, w.Body.String(), "setvar", "Expected help message to contain setvar")
}

func TestSlackCommandSetvar(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, _ := http.NewRequest("POST", "/command", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	form := url.Values{}
	form.Set("token", testToken)
	form.Set("command", "/broadway")
	form.Set("text", "setvar boing bar var1=val1")
	req.PostForm = form

	i := &instance.Instance{PlaybookID: "boing", ID: "bar", Vars: map[string]string{"var1": "val2"}}
	is := services.NewInstanceService(etcdstore.New())
	_, err := is.CreateOrUpdate(i)
	if err != nil {
		t.Log(err)
	}
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected slack command to be 200")
}

func TestSlackCommandDelete(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, _ := http.NewRequest("POST", "/command", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	form := url.Values{}
	form.Set("token", testToken)
	form.Set("command", "/broadway")
	form.Set("text", "delete helloplaybook forserver")
	req.PostForm = form

	i := &instance.Instance{PlaybookID: "helloplaybook", ID: "forserver"}
	is := services.NewInstanceService(etcdstore.New())
	_, err := is.CreateOrUpdate(i)
	if err != nil {
		t.Log(err)
	}
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected delete slack command to be 200")
}

func TestPostCommandDeployBad(t *testing.T) {
	w, _, e := helperSetupServer(testServerCfg)
	req, _ := http.NewRequest("POST", "/command", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	form := url.Values{}
	form.Set("token", testToken)
	form.Set("command", "/broadway")
	form.Set("text", "deploy foo")
	req.PostForm = form

	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected /broadway deploy foo to be 200")
	assert.Contains(t, w.Body.String(), "deploy", "Expected help message to contain deploy")
}

func TestDeployMissing(t *testing.T) {
	req, err := http.NewRequest("POST", "/deploy/missingPlaybook/missingInstance", nil)
	assert.Nil(t, err)
	req = auth(testServerCfg, req)
	w, _, e := helperSetupServer(testServerCfg)
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var errorResponse map[string]string
	log.Println(w.Body.String())
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Nil(t, err)
	assert.Contains(t, errorResponse["error"], "Not Found")
}

func TestDeleteExisting(t *testing.T) {
	ets := etcdstore.New()
	testInstance1 := &instance.Instance{
		PlaybookID: "helloplaybook",
		ID:         "TestDeleteInstance",
		Status:     instance.StatusDeployed}
	is := services.NewInstanceService(ets)
	_, err := is.CreateOrUpdate(testInstance1)
	if err != nil {
		t.Fatal(err)
	}

	req, w := testutils.DeleteRequest(
		t,
		fmt.Sprintf("/instances/%s/%s", testInstance1.PlaybookID, testInstance1.ID),
	)

	req = auth(testServerCfg, req)
	e := New(ets, testCommonCfg, testServerCfg).Handler()
	// _, _, e := helperSetupServer(testServerCfg)
	// makeRequest(s, req, w)
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected DELETE /instances to return 200")
	assert.Contains(t, w.Body.String(), "Instance successfully deleted")
}

func TestDeleteWhenNonExistantInstance(t *testing.T) {
	req, w := testutils.DeleteRequest(t, fmt.Sprintf("/%s/%s", "nonehere", "noid"))

	req = auth(testServerCfg, req)
	_, s, _ := helperSetupServer(testServerCfg)
	makeRequest(s, req, w)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected DELETE /instances to return 404 when missing instance")
}
