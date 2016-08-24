package deployment

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/glog"
	"github.com/namely/broadway/pkg/cfg"
)

const MockPlaybookContents string = `---
id: project-playbook
name: The Project 
meta:
  team: Project Devs
  email: devs@project.com
  slack: devs
vars:
  - version
  - assets_version
  - owner
messages:
  created: "help I'm stuck in a test {{ .ID }}"
  deployed: "help I'm still stuck in a test {{ .ID }}"
manifests:
  - hello
`
const MockPlaybookContentsIncomplete = `---
name: The Project 
`

const MockPlaybookContentsBadTemplate string = `---
id: project-playbook
name: The Project
messages:
  created: "help I'm stuck in a test {{ .ID }}"
  deployed: "help I'm still stuck in a test {{ .ID"
manifests:
  - hello
`
const MockPlaybookFilename = "test-playbook.yml"

var MockPlaybookBytes = []byte(MockPlaybookContents)
var MockIncompletePlaybookBytes = []byte(MockPlaybookContentsIncomplete)
var MockBadTemplatePlaybookBytes = []byte(MockPlaybookContentsBadTemplate)
var rootDir string
var playbookDir string
var mockPlaybookPath string
var testCfg = cfg.Type{
	PlaybooksPath:      "../../examples/playbooks",
	ManifestsPath:      "../../examples/manifests",
	ManifestsExtension: ".yml",
}

func SetupTestFixtures() {
	Setup(testCfg)
	// Find project root:
	cwd, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}
	rootDir = filepath.Join(cwd, "..")

	// Set path variables:
	playbookDir = filepath.Join(rootDir, "playbooks")
	mockPlaybookPath = filepath.Join(testCfg.PlaybooksPath, MockPlaybookFilename)

	// Create folders:
	err = os.MkdirAll(playbookDir, os.ModePerm)
	if err != nil {
		glog.Fatalf("Failed to create broadway/playbooks folder: %s", err)
	}

	// Write mock playbook and manifest::
	f, err := os.Create(mockPlaybookPath)
	if err != nil {
		glog.Fatalf("Failed to write mock test playbook: %s", err)
	}
	_, err = f.Write(MockPlaybookBytes)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

// TeardownTestFixtures will remove the mock files but leave the folders
func TeardownTestFixtures() {
	err := os.Remove(mockPlaybookPath)
	if err != nil {
		glog.Error(err)
	}
}

func TestMain(m *testing.M) {
	SetupTestFixtures()
	testresult := m.Run()
	TeardownTestFixtures()
	os.Exit(testresult)
}

func TestReadPlaybookFromDisk(t *testing.T) {
	playbook, err := ReadPlaybookFromDisk(mockPlaybookPath)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(playbook, MockPlaybookBytes) {
		t.Error(errors.New("Playbook read from disk differs from Playbook written to disk"))
		return
	}
}

func TestParsePlaybook(t *testing.T) {
	var err error
	ParsedPlaybook, err := ParsePlaybook(MockPlaybookBytes)
	if err != nil {
		t.Error(err)
		return
	}
	if ParsedPlaybook.Name != "The Project" {
		t.Error(errors.New("Parsed Playbook has incorrect Name field"))
		return
	}
}

func TestParsePlaybookMalformed(t *testing.T) {
	_, err := ParsePlaybook([]byte("asdf"))
	if err == nil {
		t.Error(errors.New("Parsing asdf succeeded, expected failure"))
		return
	}
}

func TestParsePlaybookIncomplete(t *testing.T) {
	_, err := ParsePlaybook(MockIncompletePlaybookBytes)
	if err != nil {
		t.Error(errors.New("Parsing well-formed, incomplete playbook failed, expected success"))
		return
	}
}

func TestValidatePlaybookPasses(t *testing.T) {
	ParsedPlaybook, _ := ParsePlaybook(MockPlaybookBytes) // already checked err in previous test

	testcases := []struct {
		scenario string
		playbook *Playbook
	}{
		{
			"Validate Valid Playbook",
			&Playbook{
				ID:        "playbook id 1",
				Name:      "playbook 1",
				Manifests: []string{"hello"},
			},
		},
		{
			"Parse And Validate Test Playbook",
			ParsedPlaybook,
		},
	}
	for _, testcase := range testcases {
		playbook := testcase.playbook
		err := playbook.Validate()
		if err != nil {
			t.Errorf("Scenario %s\nExpected: No error\nActual:\n%s", testcase.scenario, err.Error())
		}
	}
}

func TestValidatePlaybookFailures(t *testing.T) {
	InvalidPlaybook1, err := ParsePlaybook(MockBadTemplatePlaybookBytes)
	if err != nil {
		t.Error(err)
	}

	testcases := []struct {
		scenario    string
		playbook    *Playbook
		expectedErr string
	}{
		{
			"Validate Playbook Without ID",
			&Playbook{},
			"Playbook missing required ID",
		},
		{
			"Validate Playbook With Bad Template",
			InvalidPlaybook1,
			`Playbook had an invalid message template: "help I'm still stuck in a test {{ .ID"`,
		},
		{
			"Validate Playbook Without Name",
			&Playbook{
				ID: "playbook id 1",
			},
			"Playbook missing required Name",
		},
	}

	for _, testcase := range testcases {
		playbook := testcase.playbook
		err := playbook.Validate()
		if testcase.expectedErr != err.Error() {
			t.Errorf("Scenario %s\nExpected:\n%s\nActual:\n%s", testcase.scenario, testcase.expectedErr, err.Error())
		}
	}
}

func TestLoadPlaybookFolder(t *testing.T) {
	pbs, err := LoadPlaybookFolder(testCfg.PlaybooksPath)
	if err != nil {
		t.Errorf("LoadPlaybookFolder failed to load playbooks: %s\n", err)
	}
	if len(pbs) == 0 {
		t.Error("LoadPlaybookFolder failed to load")
	}
}
