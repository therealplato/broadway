package main

import (
	"github.com/namely/broadway/playbook"
)

func TestMain(m *testing.M) {
	f, err := os.Create(playbook.MockPlaybookFilename)
	if err != nil {
		log.Fatalf("Failed to write mock test playbook: %s", err)
	}
	f.Write(playbook.MockPlaybookBytes)
	f.Close()

	f, err = os.Create(playbook.MockManifestFilename)
	if err != nil {
		log.Fatalf("Failed to write mock test manifest: %s", err)
	}
	f.Close()

	testresult := m.Run()
	teardown()
	os.Exit(testresult)
}
func teardown() {
	os.Remove(playbook.MockPlaybookFilename)
	os.Remove(playbook.MockManifestFilename)
}

func Test
