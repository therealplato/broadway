package deployment

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

func setupFolder() {
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

func TestLoadManifestFolder(t *testing.T) {
	setupFolder()
	defer os.RemoveAll(tmpDir)

	ms, err := LoadManifestFolder(tmpDir, testutils.TestCfg.ManifestsExtension)
	assert.Nil(t, err)
	assert.Equal(t, len(ms), 1)
	assert.Equal(t, ms["test"].ID, "test.yml")
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
