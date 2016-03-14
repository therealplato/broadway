package playbook

import (
	"bytes"
	"errors"
	"github.com/namely/broadway/fixtures"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Switch cwd to project root before creating fixtures:
	newCwd := filepath.Join(cwd, "..")
	os.Chdir(newCwd)

	fixtures.SetupTestFixtures()
	testresult := m.Run()
	fixtures.TeardownTestFixtures()
	os.Exit(testresult)
}

func TestReadPlaybookFromDisk(t *testing.T) {
	playbook, err := ReadPlaybookFromDisk("playbooks/" + fixtures.MockPlaybookFilename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(playbook, fixtures.MockPlaybookBytes) {
		t.Error(errors.New("Playbook read from disk differs from Playbook written to disk"))
	}
}

func TestParsePlaybook(t *testing.T) {
	var err error
	ParsedPlaybook, err := ParsePlaybook(fixtures.MockPlaybookBytes)
	if err != nil {
		t.Error(err)
	}
	if ParsedPlaybook.Name != "The Project" {
		t.Error(errors.New("Parsed Playbook has incorrect Name field"))
	}
}

func TestParsePlaybookMalformed(t *testing.T) {
	_, err := ParsePlaybook([]byte("asdf"))
	if err == nil {
		t.Error(errors.New("Parsing asdf succeeded, expected failure"))
	}
}

func TestParsePlaybookIncomplete(t *testing.T) {
	_, err := ParsePlaybook(fixtures.MockIncompletePlaybookBytes)
	if err != nil {
		t.Error(errors.New("Parsing well-formed, incomplete playbook failed, expected success"))
	}
}

func TestValidatePlaybookPasses(t *testing.T) {
	ValidTask1 := Task{
		Name:      "task",
		Manifests: []string{fixtures.MockManifestFilename},
	}
	ValidTask2 := Task{
		Name:        "task",
		PodManifest: fixtures.MockManifestFilename,
	}
	ParsedPlaybook, _ := ParsePlaybook(fixtures.MockPlaybookBytes) // already checked err in previous test

	testcases := []struct {
		scenario string
		playbook Playbook
	}{
		{
			"Validate Valid Playbook",
			Playbook{
				Id:    "playbook id 1",
				Name:  "playbook 1",
				Tasks: []Task{ValidTask1, ValidTask2},
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
	InvalidTask1 := Task{
		Manifests: []string{fixtures.MockManifestFilename},
	}
	InvalidTask2 := Task{
		Name: "task",
	}

	testcases := []struct {
		scenario    string
		playbook    Playbook
		expectedErr string
	}{
		{
			"Validate Playbook Without Id",
			Playbook{},
			"Playbook missing required Id",
		},
		{
			"Validate Playbook Without Name",
			Playbook{
				Id: "playbook id 1",
			},
			"Playbook missing required Name",
		},
		{
			"Validate Playbook With Zero Tasks",
			Playbook{
				Id:    "playbook id 1",
				Name:  "playbook 1",
				Tasks: []Task{},
			},
			"Playbook requires at least 1 task",
		},
		{
			"Validate Playbook With Tasks Missing Names",
			Playbook{
				Id:    "playbook id 1",
				Name:  "playbook 1",
				Tasks: []Task{InvalidTask1},
			},
			"Task missing required Name",
		},
		{
			"Validate Playbook With Tasks Missing Manifests",
			Playbook{
				Id:    "playbook id 1",
				Name:  "playbook 1",
				Tasks: []Task{InvalidTask2},
			},
			"Task requires at least one manifest or a pod manifest",
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
	pbs, err := LoadPlaybookFolder("playbooks/")
	if err != nil {
		t.Errorf("LoadPlaybookFolder failed to load playbooks: %s\n", err)
	}
	if len(pbs) == 0 {
		t.Error("LoadPlaybookFolder failed to load mock playbook")
	}
}
func TestTaskManifestsPresentPasses(t *testing.T) {
	testcases := []struct {
		scenario string
		task     Task
	}{
		{
			"Task With Existing Manifests",
			Task{
				Name:      "task 2",
				Manifests: []string{fixtures.MockManifestFilename},
			},
		},
		{
			"Task With Only Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: fixtures.MockManifestFilename,
			},
		},
		{
			"Task With Both Manifests And Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: fixtures.MockManifestFilename,
				Manifests:   []string{fixtures.MockManifestFilename},
			},
		},
	}
	for _, testcase := range testcases {
		task := testcase.task
		err := task.ManifestsPresent()
		if err != nil {
			t.Errorf("Scenario %s\nExpected: Success\nActual:\n%s", testcase.scenario, err.Error())
		}
	}
}
func TestTaskManifestsPresentFailures(t *testing.T) {
	testcases := []struct {
		scenario string
		task     Task
	}{
		{
			"Task With Missing Manifests",
			Task{
				Name:      "task 1",
				Manifests: []string{"pod0"},
			},
		},
	}

	for _, testcase := range testcases {
		task := testcase.task
		err := task.ManifestsPresent()
		if err == nil {
			t.Errorf("Scenario %s\nExpected: File does not exist\nActual: Success%s", testcase.scenario)
		} else if !os.IsNotExist(err) { // it was the wrong error!
			t.Errorf("Scenario %s\nExpected: File does not exist\nActual:\n%s", testcase.scenario, err.Error())
		}
	}
}
