// Test fixtures shared throughout Broadway tests
// Do not import this except from tests
package fixtures

import (
	"log"
	"os"
	"path/filepath"
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

const MockPlaybookContentsIncomplete = `---
name: The Project 
`

var MockPlaybookBytes = []byte(MockPlaybookContents)
var MockIncompletePlaybookBytes = []byte(MockPlaybookContentsIncomplete)

func SetupTestFixtures() {
	// Ensure we are in project root:
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.Base(cwd) != "broadway" {
		log.Fatalf("Failed to setup test fixtures; expected cwd 'broadway/', actual cwd %s", cwd)
	}
	// Ensure playbooks and manifests folders to write mock data
	pDir := filepath.Join(cwd, "playbooks")
	mDir := filepath.Join(cwd, "manifests")
	err = os.MkdirAll(pDir, os.ModePerm)
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
}

func TeardownTestFixtures() {
	os.Remove(MockPlaybookFilename)
	os.Remove(MockManifestFilename)
}
