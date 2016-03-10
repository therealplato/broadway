package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

const PlaybookFilename = "test-playbook.yml"
const ManifestFilename = "test-manifest.yml"
const PlaybookContents string = `---
name: The Project 
meta:
  team: Project Devs
  email: devs@project.com
  slack: devs
vars:
  - version
  - assets_version
  - owner
tasks:
  - name: Deploy Postgres
    manifests:
      - test-manifest.yml
      - test-manifest.yml
  - name: Deploy Redis
    manifests:
      - test-manifest.yml
      - test-manifest.yml
  - name: Database Setup
    pod_manifest: test-manifest.yml
    wait_for:
      - success
    when: new_deployment
  - name: Database Migration
    pod_manifest: test-manifest.yml
    wait_for:
      - success
  - name: Deploy Project
    manifests:
      - test-manifest.yml
      - test-manifest.yml
      - test-manifest.yml
`

var PlaybookBytes = []byte(PlaybookContents)

const PlaybookContentsIncomplete = `---
name: The Project 
`

var IncompletePlaybookBytes = []byte(PlaybookContentsIncomplete)

func TestMain(m *testing.M) {
	f, _ := os.Create(PlaybookFilename)
	f.Write(PlaybookBytes)
	f.Close()

	f, _ = os.Create(ManifestFilename)
	f.Close()

	testresult := m.Run()
	teardown()
	os.Exit(testresult)
}
func teardown() {
	os.Remove(PlaybookFilename)
	os.Remove(ManifestFilename)
}

func TestReadPlaybookFromDisk(t *testing.T) {
	playbook, err := ReadPlaybookFromDisk(PlaybookFilename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(playbook, PlaybookBytes) {
		t.Error(errors.New("Playbook read from disk differs from Playbook written to disk"))
	}
}

func TestParsePlaybook(t *testing.T) {
	var err error
	ParsedPlaybook, err := ParsePlaybook(PlaybookBytes)
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
	_, err := ParsePlaybook(IncompletePlaybookBytes)
	if err != nil {
		t.Error(errors.New("Parsing well-formed, incomplete playbook failed, expected success"))
	}
}

func TestValidatePlaybook(t *testing.T) {
	ValidTask := Task{
		Name:      "task",
		Manifests: []string{ManifestFilename},
	}
	InvalidTask1 := Task{
		Manifests: []string{ManifestFilename},
	}
	InvalidTask2 := Task{
		Name: "task",
	}

	testcases := []struct {
		scenario    string
		playbook    Playbook
		errExpected bool
	}{
		{
			"Validate Valid Playbook",
			Playbook{
				Name:  "playbook 1",
				Tasks: []Task{ValidTask},
			},
			false,
		},
		{
			"Validate Empty Playbook",
			Playbook{},
			true,
		},
		{
			"Validate Playbook With Tasks Missing Names",
			Playbook{
				Name:  "playbook 1",
				Tasks: []Task{InvalidTask1},
			},
			true,
		},
		{
			"Validate Playbook With Tasks Missing Manifests",
			Playbook{
				Name:  "playbook 1",
				Tasks: []Task{InvalidTask2},
			},
			true,
		},
	}

	for _, testcase := range testcases {
		playbook := testcase.playbook
		err := playbook.Validate()
		if testcase.errExpected && err != nil {
			continue
		} else if testcase.errExpected && err == nil {
			t.Errorf("Scenario %s: Succeeded, expected failure", testcase.scenario)
		} else if err != nil {
			t.Errorf("Scenario %s: Failed with %s, expected success", testcase.scenario, err)
		}
	}
}

func TestTaskManifestsPresent(t *testing.T) {
	testcases := []struct {
		scenario    string
		task        Task
		errExpected bool
	}{
		{
			"Task With Missing Manifests",
			Task{
				Name:      "task 1",
				Manifests: []string{"pod0"},
			},
			true,
		},
		{
			"Task With Existing Manifests",
			Task{
				Name:      "task 2",
				Manifests: []string{ManifestFilename},
			},
			false,
		},
		{
			"Task With Only Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: ManifestFilename,
			},
			false,
		},
	}

	for _, testcase := range testcases {
		task := testcase.task
		err := task.ManifestsPresent()
		if testcase.errExpected {
			if !os.IsNotExist(err) { // it was the wrong error!
				t.Errorf("Scenario %s: Got %s, expected 'not found'", testcase.scenario, err)
			}
		} else if err != nil {
			t.Errorf("Scenario %s: Got %s, expected success", testcase.scenario, err)
		}
	}
}
