package main

import (
	"os"
	"testing"
)

const ManifestFilename = "test-manifest.yml"
const Manifest string = `---
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
	f, _ := os.Create(ManifestFilename)
	f.Write([]byte(Manifest))
	f.Close()
	tres := m.Run()
	teardown()
	os.Exit(tres)
}
func teardown() {
	os.Remove(ManifestFilename)
}

func TestValidateEmptyManifests(t *testing.T) {

	TestTaskWithoutManifest := Task{
		Name: ManifestFilename,
	}
	err := TestTaskWithoutManifest.ValidateManifests()
	if err != nil {
		t.Error(err)
	}
}

func TestValidateManifests(t *testing.T) {

	TestTaskWithManifest := Task{
		Name:      ManifestFilename,
		Manifests: []string{ManifestFilename},
	}
	err := TestTaskWithManifest.ValidateManifests()
	if err != nil {
		t.Error(err)
	}
}
