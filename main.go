package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Meta struct {
	Team  string `yaml:"team"`
	Email string `yaml:"email"`
	Slack string `yaml:"slack"`
}

type Task struct {
	Name        string   `yaml:"name"`
	Manifests   []string `yaml:"manifests,omitempty"`
	PodManifest []string `yaml:"pod_manifest,omitempty"`
	WaitFor     []string `yaml:"wait_for,omitempty"`
	When        string   `yaml:"when,omitempty"`
}

type Playbook struct {
	Name  string   `yaml:"name"`
	Meta  Meta     `yaml:"meta"`
	Vars  []string `yaml:"vars"`
	Tasks []Task   `yaml:"tasks"`
}

func (t Task) ManifestsPresent() error {
	for _, name := range t.Manifests {
		if _, err := os.Stat(name); err != nil {
			return err
		}
	}
	return nil
}

func (p Playbook) ValidateTasks() error {
	for _, task := range p.Tasks {
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

func main() {
	args := os.Args
	yamlFileDescriptor := args[1:][0]

	playbookBytes, err := ReadPlaybookFromDisk(yamlFileDescriptor)
	if err != nil {
		log.Fatal(err)
	}

	playbook, err := ParsePlaybook(playbookBytes)
	if err != nil {
		log.Fatal(err)
	}

	if err := playbook.ValidateTasks(); err != nil {
		log.Fatalf("Task validation failed: %s", err)
	}

	fmt.Printf("%+v", playbook)
}
