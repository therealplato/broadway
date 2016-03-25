package services

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmpDir string

func Setup() {
	var err error
	tmpDir, err = ioutil.TempDir("", "manifest_test_")
	if err != nil {
		log.Fatal(err)
	}
	content := []byte(testManifest)
	tmpFN := filepath.Join(tmpDir, "test.yml")
	if err := ioutil.WriteFile(tmpFN, content, 0666); err != nil {
		log.Fatal(err)
	}
}
func Teardown() {
	defer os.RemoveAll(tmpDir)
}

func TestMain(m *testing.M) {
	Setup()
	testresult := m.Run()
	Teardown()
	os.Exit(testresult)
}

func TestNewManifestService(t *testing.T) {
	ms := NewManifestService()
	assert.Equal(t, "./manifests", ms.rootFolder)
}

func TestRead(t *testing.T) {
	ms := NewManifestService()
	ms.rootFolder = tmpDir
	assert.Equal(t, tmpDir, ms.rootFolder)
	contents, err := ms.Read("test")
	assert.Contains(t, contents, "ReplicationController")
	assert.Nil(t, err)
}

const testManifest = `apiVersion: v1
kind: ReplicationController
metadata:
  name: test
spec:
  replicas: 1
  selector:
    name: redis
  template:
    metadata:
      labels:
        name: redis
    spec:
      containers:
      - name: redis
        image: kubernetes/redis:v1
`
