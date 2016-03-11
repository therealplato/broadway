package playbook

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const MockPlaybookFilename = "playbooks/test-playbook.yml"
const MockManifestFilename = "manifests/test-manifest.yml"
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

var MockPlaybookBytes = []byte(MockPlaybookContents)

const MockPlaybookContentsIncomplete = `---
name: The Project 
`

var MockIncompletePlaybookBytes = []byte(MockPlaybookContentsIncomplete)

func TestMain(m *testing.M) {
	// Switch to project root. The mock files are relative to there.
	cwd, _ := os.Getwd()
	newCwd := filepath.Join(cwd, "..")
	os.Chdir(newCwd)
	fmt.Println(os.Getwd())

	// Ensure playbooks and manifests folders to write mock data
	pDir := filepath.Join(newCwd, "playbooks")
	mDir := filepath.Join(newCwd, "manifests")
	err := os.MkdirAll(pDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create broadway/playbooks folder: %s", err)
	}
	err = os.MkdirAll(mDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create broadway/manifests folder: %s", err)
	}

	f, err := os.Create(MockPlaybookFilename)
	if err != nil {
		log.Fatalf("Failed to write mock test playbook: %s", err)
	}
	f.Write(MockPlaybookBytes)
	f.Close()

	f, err = os.Create(MockManifestFilename)
	if err != nil {
		log.Fatalf("Failed to write mock test manifest: %s", err)
	}
	f.Close()

	testresult := m.Run()
	teardown()
	os.Exit(testresult)
}
func teardown() {
	os.Remove(MockPlaybookFilename)
	os.Remove(MockManifestFilename)
}

func TestReadPlaybookFromDisk(t *testing.T) {
	playbook, err := ReadPlaybookFromDisk(MockPlaybookFilename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(playbook, MockPlaybookBytes) {
		t.Error(errors.New("Playbook read from disk differs from Playbook written to disk"))
	}
}

func TestParsePlaybook(t *testing.T) {
	var err error
	ParsedPlaybook, err := ParsePlaybook(MockPlaybookBytes)
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
	_, err := ParsePlaybook(MockIncompletePlaybookBytes)
	if err != nil {
		t.Error(errors.New("Parsing well-formed, incomplete playbook failed, expected success"))
	}
}

func TestValidatePlaybookPasses(t *testing.T) {
	ValidTask1 := Task{
		Name:      "task",
		Manifests: []string{MockManifestFilename},
	}
	ValidTask2 := Task{
		Name:        "task",
		PodManifest: MockManifestFilename,
	}
	ParsedPlaybook, _ := ParsePlaybook(MockPlaybookBytes) // already checked err in previous test

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
		Manifests: []string{MockManifestFilename},
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

func TestTaskManifestsPresentPasses(t *testing.T) {
	testcases := []struct {
		scenario string
		task     Task
	}{
		{
			"Task With Existing Manifests",
			Task{
				Name:      "task 2",
				Manifests: []string{MockManifestFilename},
			},
		},
		{
			"Task With Only Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: MockManifestFilename,
			},
		},
		{
			"Task With Both Manifests And Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: MockManifestFilename,
				Manifests:   []string{MockManifestFilename},
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
