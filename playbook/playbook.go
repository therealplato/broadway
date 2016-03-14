package playbook

import (
	"errors"
	"fmt"
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

var ManifestRoot = "manifests/"
var ManifestExtension = ".yml"

func SetManifestRoot(newRoot string) error {
	if _, err := os.Stat(newRoot); err != nil {
		return err
	}
	ManifestRoot = newRoot
	return nil
}

func (t Task) ManifestsPresent() error {
	for _, name := range t.Manifests {
		filename := name + ManifestExtension
		path := filepath.Join(ManifestRoot, filename)
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	if len(t.PodManifest) > 0 {
		filename := t.PodManifest + ManifestExtension
		path := filepath.Join(ManifestRoot, filename)
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
func LoadPlaybookFolder(dir string) ([]Playbook, error) {
	var AllPlaybooks []Playbook
	paths, err := filepath.Glob(dir + "/*")
	if err != nil {
		return AllPlaybooks, err
	}
	for _, path := range paths {
		playbookBytes, err := ReadPlaybookFromDisk(path)
		if err != nil {
			fmt.Printf("Warning: Failed to read %s\n", path)
			continue
		}
		parsed, err := ParsePlaybook(playbookBytes)
		if err != nil {
			fmt.Printf("Warning: Failed to parse %s\n", path)
			continue
		}
		err = parsed.Validate()
		if err != nil {
			fmt.Printf("Warning: Playbook %s invalid: %s\n", path, err)
			continue
		}
		AllPlaybooks = append(AllPlaybooks, parsed)
	}
	return AllPlaybooks, nil
}
