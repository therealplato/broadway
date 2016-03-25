package services

import (
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
)

// ManifestService mediates between things that use manifests and manifest
// implementations
type ManifestService struct {
	rootFolder string
}

// NewManifestService instantiates a ManifestService with a default rootFolder
func NewManifestService() *ManifestService {
	return &ManifestService{
		rootFolder: "./manifests",
	}
}

// Read loads the contents of `name` from disk, looking in the
// ManifestService.rootFolder
func (ms *ManifestService) Read(name string) (string, error) {
	return "dummy", nil
}

// LoadTask iterates through the podManifest and Manifests of a Task and returns
// Manifest objects
func (ms *ManifestService) LoadTask(t playbook.Task) (*manifest.Manifest, []manifest.Manifest, error) {
	pm := &manifest.Manifest{}
	var mm []manifest.Manifest
	if len(t.PodManifest) != 0 {
		pm, err := ms.Load(t.PodManifest)
		if err != nil {
			return pm, mm, err
		}
	}
	for _, name := range t.Manifests {
		m, err := ms.Load(name)
		if err != nil {
			return pm, mm, err
		}
		mm = append(mm, *m)
	}

	return pm, mm, nil
}

// Load takes a filename, reads the file and generates a Manifest object
func (ms *ManifestService) Load(name string) (*manifest.Manifest, error) {
	var mPath = filepath.Join(ms.rootFolder, name)
	bytes, err := ioutil.ReadFile(mPath)
	if err != nil {
		return &manifest.Manifest{}, err
	}
	m, err := ms.New(name, string(bytes))
	if err != nil {
		return &manifest.Manifest{}, err
	}
	return m, nil
}

// New creates a new Manifest object and parses (but does not execute) the template
func (ms *ManifestService) New(id, content string) (*manifest.Manifest, error) {
	t, err := template.New(id).Parse(content)
	if err != nil {
		return nil, err
	}

	return &manifest.Manifest{ID: id, Template: t}, nil
}
