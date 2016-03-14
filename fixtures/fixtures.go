// Test fixtures shared throughout Broadway tests
// Do not import this except from tests
package fixtures

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
