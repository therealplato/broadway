package main

import (
	"fmt"
	"github.com/namely/broadway/playbook"
	//"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println(playbook.MockPlaybookFilename)
	/*
			fmt.Println(playbook.MockPlaybookFilename)
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

		cwd, _ := os.Getwd()
		playbook.SetupTestFixtures()
		testresult := m.Run()
		os.Exit(testresult)
		playbook.TeardownTestFixtures()
	*/
}

/*
func teardown() {
	os.Remove(playbook.MockPlaybookFilename)
	os.Remove(playbook.MockManifestFilename)
}
*/

func TestLoadPlaybookFolder(t *testing.T) {
}
