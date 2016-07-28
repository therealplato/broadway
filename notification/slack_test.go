package notification

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/namely/broadway/cfg"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	requestBody := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("No request Body received")
		}

		requestBody = string(contents)
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()

	message := &Message{
		Attachments: []Attachment{{
			Text: "successful",
		}},
		cfg: cfg.ServerCfgType{
			SlackWebhook: ts.URL,
		},
	}
	err := message.Send()
	assert.Nil(t, err)
	assert.Contains(t, requestBody, "successful")
}
