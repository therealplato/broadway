package services

import (
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
)

// ManifestService mediates between things that use manifests and manifest
// implementations
type ManifestService struct {
	rootFolder string
	extension  string
}

// NewManifestService instantiates a ManifestService with a default rootFolder
func NewManifestService() *ManifestService {
	return &ManifestService{
		rootFolder: "./manifests",
		extension:  ".yml",
	}
}

// Read loads the contents of `name` from disk, looking in the
// ManifestService.rootFolder
func (ms *ManifestService) Read(name string) (string, error) {
	name = name + ms.extension
	path := filepath.Join(ms.rootFolder, name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// LoadTask iterates through the podManifest and Manifests of a Task and returns
// Manifest objects
func (ms *ManifestService) LoadTask(t playbook.Task) (*manifest.Manifest, []manifest.Manifest, error) {
	pm := &manifest.Manifest{}
	var mm []manifest.Manifest
	if len(t.PodManifest) != 0 {
		var err error
		pm, err = ms.Load(t.PodManifest)
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
	mString, err := ms.Read(name)
	if err != nil {
		return &manifest.Manifest{}, err
	}
	m, err := ms.New(name, mString)
	if err != nil {
		return &manifest.Manifest{}, err
	}
	return m, nil
}

// New creates a new Manifest object and parses (but does not execute) the template
func (ms *ManifestService) New(id, content string) (*manifest.Manifest, error) {
	return manifest.New(id, content)
}

// LoadManifestFolder returns a map of name:Manifest for each file in
// ms.rootFolder
func (ms *ManifestService) LoadManifestFolder() (map[string]*manifest.Manifest, error) {
	var mm = make(map[string]*manifest.Manifest)
	extRX, err := regexp.Compile(`\.[^.]+$`)
	paths, err := filepath.Glob(filepath.Join(ms.rootFolder, "*"))
	if err != nil {
		return mm, err
	}
	for _, p := range paths {
		filename := filepath.Base(p)
		name := extRX.ReplaceAllString(filename, "") // remove extension
		m, err := ms.Load(name)
		if err != nil {
			return mm, err
		}
		mm[name] = m
	}
	return mm, nil
}
