package deployment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/golang/glog"
	"github.com/namely/broadway/pkg/cfg"

	"gopkg.in/yaml.v2"
)

// Meta contains optional metadata keys associated with this playbook
type Meta struct {
	Team  string `yaml:"team"`
	Email string `yaml:"email"`
	Slack string `yaml:"slack"`
}

// Playbook configures a set of tasks to be automated
type Playbook struct {
	ID        string            `yaml:"id"`
	Name      string            `yaml:"name"`
	Meta      Meta              `yaml:"meta"`
	Vars      []string          `yaml:"vars"`
	Manifests []string          `yaml:"manifests"`
	Messages  map[string]string `yaml:"messages"`
}

// AllPlaybooks is a map of playbook id's to playbooks
var AllPlaybooks map[string]*Playbook
var playbooksPath string
var manifestsPath string
var manifestsExtension string

// SetupPlaybook configures playbook with an injected configuration
func SetupPlaybook(cfg cfg.Type) {
	playbooksPath = cfg.PlaybooksPath
	manifestsPath = cfg.ManifestsPath
	manifestsExtension = cfg.ManifestsExtension
	var err error
	AllPlaybooks, err = LoadPlaybookFolder(cfg.PlaybooksPath)
	if err != nil {
		glog.Fatal(err)
	}
}

// Validate checks for ID, Name, and Tasks on a playbook
func (p *Playbook) Validate() error {
	if len(p.ID) == 0 {
		return errors.New("Playbook missing required ID")
	}
	if len(p.Name) == 0 {
		return errors.New("Playbook missing required Name")
	}
	if len(p.Manifests) == 0 {
		return errors.New("Playbook requires at least 1 manifest")
	}
	for key, value := range p.Messages {
		_, err := template.New(key).Parse(value)
		if err != nil {
			return fmt.Errorf("Playbook had an invalid message template: \"%s\"", value)
		}
	}
	return p.ValidateManifests()
}

// ValidateManifests checks manifests
func (p *Playbook) ValidateManifests() error {
	for _, name := range p.Manifests {
		filename := name + manifestsExtension
		path := filepath.Join(manifestsPath, filename)
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

// ParsePlaybook unmarshalls a YAML byte sequence into a Playbook struct
func ParsePlaybook(playbook []byte) (*Playbook, error) {
	var p Playbook
	err := yaml.Unmarshal(playbook, &p)
	return &p, err
}

// ReadPlaybookFromDisk takes a filename and returns a byte array. Alias for
// ioutil.ReadFile
func ReadPlaybookFromDisk(fd string) ([]byte, error) {
	return ioutil.ReadFile(fd)
}

// LoadPlaybookFolder takes a directory and attempts to parse every file in that
// directory into a Playbook struct
func LoadPlaybookFolder(dir string) (map[string]*Playbook, error) {
	var playbooks = make(map[string]*Playbook)
	paths, err := filepath.Glob(dir + "/*")
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, errors.New("Found zero files in directory " + dir)
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
		playbooks[parsed.ID] = parsed
	}
	return playbooks, nil
}
