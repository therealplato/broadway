package manifest

import (
	"bytes"
	"text/template"
)

// Manifest represents a kubernetes manifest file
// Filename is used as identifier in the current implementation.
type Manifest struct {
	ID       string
	Template *template.Template
}

// Execute executes template with variables
func (m *Manifest) Execute(vars map[string]string) string {
	var b bytes.Buffer
	err := m.Template.Execute(&b, vars)
	if err != nil {
		return ""
	}
	return b.String()
}
