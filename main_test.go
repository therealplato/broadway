package main

import (
	"github.com/namely/broadway/fixtures"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fixtures.SetupTestFixtures()
	testresult := m.Run()
	os.Exit(testresult)
	fixtures.TeardownTestFixtures()
}

func TestLoadPlaybookFolder(t *testing.T) {
	pbs, err := LoadPlaybookFolder("playbooks/")
	if err != nil {
		t.Errorf("LoadPlaybookFolder failed to load playbooks: %s\n", err)
	}
	if len(pbs) == 0 {
		t.Error("LoadPlaybookFolder failed to load mock playbook")
	}
}
