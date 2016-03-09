package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

var ParsedPlaybook Playbook

const PlaybookFilename = "test-manifest.yml"
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
      - postgres-rc
      - postgres-service
  - name: Deploy Redis
    manifests:
      - redis-rc
      - redis-service
  - name: Database Setup
    pod_manifest:
      - createdb-pod
    wait_for:
      - success
    when: new_deployment
  - name: Database Migration
    pod_manifest:
      - migration-pod
    wait_for:
      - success
  - name: Deploy Project
    manifests:
      - web-rc
      - web-service
      - sidekiq-rc
`

var PlaybookBytes = []byte(PlaybookContents)

func TestMain(m *testing.M) {
	f, _ := os.Create(PlaybookFilename)
	f.Write(PlaybookBytes)
	//[]byte(PlaybookContents))
	f.Close()
	testresult := m.Run()
	teardown()
	os.Exit(testresult)
}
func teardown() {
	os.Remove(PlaybookFilename)
}

func TestReadPlaybookFromDisk(t *testing.T) {
	playbook, err := ReadPlaybookFromDisk(PlaybookFilename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(playbook, PlaybookBytes) {
		//[]byte(PlaybookContents) {
		t.Error(errors.New("Playbook read from disk differs from Playbook written to disk"))
	}
}

func TestParsePlaybook(t *testing.T) {
	var err error
	ParsedPlaybook, err = ParsePlaybook(PlaybookBytes)
	if err != nil {
		t.Error(err)
	}
	// Todo: Pass in malformed yaml
	// Todo: Parse in well-formed yaml missing required fields
}

func TestValidatePlaybook(t *testing.T) {
	if ParsedPlaybook.Name != "The Project" {
		t.Error(errors.New("Parsed Playbook has incorrect Name field"))
	}
	// Todo: Write and test ValidatePlaybook
}

func TestTaskManifestsPresent(t *testing.T) {
	testcases := []struct {
		scenario    string
		task        Task
		errExpected bool
	}{
		{
			"Validate Task Without Manifests",
			Task{
				Name: "task 0",
			},
			false,
		},
		{
			"Validate Task With Missing Manifests",
			Task{
				Name:      "task 1",
				Manifests: []string{"pod0"},
			},
			true,
		},
		{
			"Validate Task With Existing Manifests",
			Task{
				Name:      "task 2",
				Manifests: []string{PlaybookFilename},
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
