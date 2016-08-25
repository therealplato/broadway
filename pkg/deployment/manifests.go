package deployment

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var extExp = regexp.MustCompile(`\.[^.]+$`)

// LoadManifestFolder returns a map of name:Manifest for each file in root
func LoadManifestFolder(root, ext string) (map[string]*Manifest, error) {
	var mm = make(map[string]*Manifest)
	paths, err := filepath.Glob(filepath.Join(root, "*"))
	if err != nil {
		return mm, err
	}
	for _, p := range paths {
		filename := filepath.Base(p)
		name := extExp.ReplaceAllString(filename, "") // remove ext
		m, err := load(root, name+ext)
		if err != nil {
			return mm, err
		}
		mm[name] = m
	}
	return mm, nil
}

func load(root, name string) (*Manifest, error) {
	path := filepath.Join(root, name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	if err != nil {
		return &Manifest{}, err
	}
	m, err := NewManifest(name, content)
	if err != nil {
		return &Manifest{}, err
	}
	return m, nil
}
