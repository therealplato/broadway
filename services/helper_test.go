package services

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/namely/broadway/testutils"
)

// ServicesTestCfg is a config created for services tests that can be safely modified
var ServicesTestCfg = testutils.TestCfg

type notificationTestHelper struct {
	requestBody string
	ts          *httptest.Server
}

func newNotificationTestHelper() *notificationTestHelper {
	n := &notificationTestHelper{
		requestBody: "",
	}
	n.ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("No request Body received")
		}

		n.requestBody = string(contents)
	}))
	ServicesTestCfg.SlackWebhook = n.ts.URL
	return n
}

func (n *notificationTestHelper) Close() {
	n.ts.Close()
}
