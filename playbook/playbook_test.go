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

const MockPlaybookFilename = "test-playbook.yml"
const MockManifestFilename = "test-manifest.yml"
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
      - test-manifest
      - test-manifest
  - name: Deploy Redis
    manifests:
      - test-manifest
      - test-manifest
  - name: Database Setup
    pod_manifest: test-manifest
    wait_for:
      - success
    when: new_deployment
  - name: Database Migration
    pod_manifest: test-manifest
    wait_for:
      - success
  - name: Deploy Project
    manifests:
      - test-manifest
      - test-manifest
      - test-manifest
`
const MockPlaybookContentsIncomplete = `---
name: The Project 
`

var MockPlaybookBytes = []byte(MockPlaybookContents)
var MockIncompletePlaybookBytes = []byte(MockPlaybookContentsIncomplete)
var rootPath string

// SetupTestFixtures will write broadway/playbooks/test-playbook.yml
// and broadway/playbooks/test-manifest.yml, creating folders if necessary
func SetupTestFixtures() {
	// Ensure we are in project root:
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.Base(cwd) != "broadway" {
		log.Fatalf("Failed to setup test fixtures; expected cwd 'broadway/', actual cwd %s", cwd)
	}
	rootPath = cwd
	// Write mock playbook:
	pDir := filepath.Join(rootPath, "playbooks")
	err = os.MkdirAll(pDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create broadway/playbooks folder: %s", err)
	}
	err = os.Chdir(pDir)
	if err != nil {
		log.Fatalf("Failed to cd to broadway/playbooks folder: %s", err)
	}
	f, err := os.Create(MockPlaybookFilename)
	if err != nil {
		log.Fatalf("Failed to write mock test playbook: %s", err)
	}
	f.Write(MockPlaybookBytes)
	f.Close()

	// Write mock manifest:
	mDir := filepath.Join(rootPath, "manifests")
	err = os.MkdirAll(mDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create broadway/manifests folder: %s", err)
	}
	err = os.Chdir(mDir)
	if err != nil {
		log.Fatalf("Failed to cd to broadway/manifests folder: %s", err)
	}

	f, err = os.Create(MockManifestFilename)
	if err != nil {
		log.Fatalf("Failed to write mock test manifest: %s", err)
	}
	f.Close()

	err = os.Chdir(rootPath)
	if err != nil {
		log.Fatalf("Failed to cd to broadway/ folder: %s", err)
	}
}

// TeardownTestFixtures will remove the mock files but leave the folders
func TeardownTestFixtures() {
	pPath := filepath.Join(rootPath, "playbooks", MockPlaybookFilename)
	mPath := filepath.Join(rootPath, "manifests", MockManifestFilename)
	err1 := os.Remove(pPath)
	err2 := os.Remove(mPath)
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
	}
}

func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Switch cwd to project root before creating fixtures:
	newCwd := filepath.Join(cwd, "..")
	os.Chdir(newCwd)

	SetupTestFixtures()
	testresult := m.Run()
	TeardownTestFixtures()
	os.Exit(testresult)
}

func TestReadPlaybookFromDisk(t *testing.T) {
	playbook, err := ReadPlaybookFromDisk("playbooks/" + MockPlaybookFilename)
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
		Manifests: []string{"test-manifest"},
	}
	ValidTask2 := Task{
		Name:        "task",
		PodManifest: "test-manifest",
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
				Manifests: []string{"test-manifest"},
			},
		},
		{
			"Task With Only Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: "test-manifest",
			},
		},
		{
			"Task With Both Manifests And Pod Manifest",
			Task{
				Name:        "task 2",
				PodManifest: "test-manifest",
				Manifests:   []string{"test-manifest"},
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
