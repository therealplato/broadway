package manifest

import (
	"bytes"
	"text/template"
)

// Manifest represents a kubernetes manifest file
type Manifest struct {
	Name     string
	template *template.Template
}

// New creates a new Manifest object and parses the template
func New(name, content string) (*Manifest, error) {
	t, err := template.New(name).Parse(content)
	if err != nil {
		return nil, err
	}

	return &Manifest{Name: name, template: t}, nil
}

// Execute runs teamplate with variables
func (m *Manifest) Execute(vars map[string]string) string {
	var b bytes.Buffer
	err := m.template.Execute(&b, vars)
	if err != nil {
		return ""
	}
	return b.String()
}
