package playbook

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Meta struct {
	Team  string `yaml:"team"`
	Email string `yaml:"email"`
	Slack string `yaml:"slack"`
}

type Task struct {
	Name        string   `yaml:"name"`
	Manifests   []string `yaml:"manifests,omitempty"`
	PodManifest string   `yaml:"pod_manifest,omitempty"`
	WaitFor     []string `yaml:"wait_for,omitempty"`
	When        string   `yaml:"when,omitempty"`
}

type Playbook struct {
	Id    string   `yaml:"id"`
	Name  string   `yaml:"name"`
	Meta  Meta     `yaml:"meta"`
	Vars  []string `yaml:"vars"`
	Tasks []Task   `yaml:"tasks"`
}

var PlaybookRoot = "playbooks/"
var ManifestRoot = "manifests/"

func SetPlaybookRoot(newRoot string) error {
	if _, err := os.Stat(newRoot); err != nil {
		return err
	}
	PlaybookRoot = newRoot
	return nil
}
func SetManifestRoot(newRoot string) error {
	if _, err := os.Stat(newRoot); err != nil {
		return err
	}
	ManifestRoot = newRoot
	return nil
}

func (t Task) ManifestsPresent() error {
	for _, name := range t.Manifests {
		path := filepath.Join(ManifestRoot, name)
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	if len(t.PodManifest) > 0 {
		path := filepath.Join(ManifestRoot, t.PodManifest)
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

func (p Playbook) Validate() error {
	if len(p.Id) == 0 {
		return errors.New("Playbook missing required Id")
	}
	if len(p.Name) == 0 {
		return errors.New("Playbook missing required Name")
	}
	if len(p.Tasks) == 0 {
		return errors.New("Playbook requires at least 1 task")
	}
	return p.ValidateTasks()
}

func (p Playbook) ValidateTasks() error {
	for _, task := range p.Tasks {
		if len(task.Name) == 0 {
			return errors.New("Task missing required Name")
		}
		if len(task.Manifests) == 0 && len(task.PodManifest) == 0 {
			return errors.New("Task requires at least one manifest or a pod manifest")
		}
		if err := task.ManifestsPresent(); err != nil {
			return err
		}
	}
	return nil
}

func ParsePlaybook(playbook []byte) (Playbook, error) {
	var p Playbook
	err := yaml.Unmarshal(playbook, &p)
	return p, err
}

func ReadPlaybookFromDisk(fd string) ([]byte, error) {
	return ioutil.ReadFile(fd)
}
