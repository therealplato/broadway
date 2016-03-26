package manifest

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	id := "test"
	tp, err := template.New(id).Parse(`{{ .test }}`)
	assert.Nil(t, err)
	m := &Manifest{ID: id, Template: tp}

	out := m.Execute(map[string]string{"test": "hello!"})

	assert.Equal(t, "hello!", out)
}
