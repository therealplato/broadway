package main

import (
	"os"
	"testing"
)

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

func TestPass(t *testing.T) {
	//t.Succeed()
}
func TestFail(t *testing.T) {
	//t.Fail()
}

func TestMain(m *testing.M) {
	f, _ := os.Create(PlaybookFilename)
	f.Write([]byte(PlaybookContents))
	f.Close()
	tres := m.Run()
	teardown()
	os.Exit(tres)
}
func teardown() {
	os.Remove(PlaybookFilename)
}

func TestValidateManifests(t *testing.T) {
	testcases := []struct {
		scenario string
		task     Task
		//expectedErr error
		errExpected bool
	}{
		{
			"TestTaskWithoutManifests",
			Task{
				Name: "task 0",
			},
			false,
		},
		{
			"TestTaskWithMissingManifests",
			Task{
				Name:      "task 1",
				Manifests: []string{"pod0"},
			},
			true,
		},
		{
			"TestTaskWithExistingManifests",
			Task{
				Name:      "task 2",
				Manifests: []string{PlaybookFilename},
			},
			false,
		},
	}

	for _, testcase := range testcases {
		task := testcase.task
		err := task.ValidateManifests()
		if testcase.errExpected {
			if !os.IsNotExist(err) { // it was the wrong error!
				t.Errorf("Scenario %s: Got %s, expected 'not found'", testcase.scenario, err)
			}
		} else if err != nil {
			t.Errorf("Scenario %s: Got %s, expected success", testcase.scenario, err)
		}
	}
}
