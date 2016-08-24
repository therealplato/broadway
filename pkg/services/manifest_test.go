package services

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/namely/broadway/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

var tmpDir string

func setup() {
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
	setup()
	testresult := m.Run()
	Teardown()
	os.Exit(testresult)
}

func TestNewManifestService(t *testing.T) {
	ms := NewManifestService(testutils.TestCfg)
	assert.Equal(t, testutils.TestCfg.ManifestsPath, ms.rootFolder)
}

func TestRead(t *testing.T) {
	ms := NewManifestService(testutils.TestCfg)
	ms.rootFolder = tmpDir
	assert.Equal(t, tmpDir, ms.rootFolder)
	contents, err := ms.Read("test")
	assert.Contains(t, contents, "ReplicationController")
	assert.Nil(t, err)

	_, err = ms.Read("missing")
	assert.NotNil(t, err)
}

func TestLoad(t *testing.T) {
	ms := NewManifestService(testutils.TestCfg)
	ms.rootFolder = tmpDir
	m, err := ms.Load("test")
	assert.Nil(t, err)
	assert.Equal(t, m.ID, "test")

	m, err = ms.Load("missing")
	assert.NotNil(t, err)
}

func TestNew(t *testing.T) {
	ms := NewManifestService(testutils.TestCfg)
	ms.rootFolder = tmpDir

	m, err := ms.New("testId", "testContent")
	assert.Equal(t, m.ID, "testId")
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
