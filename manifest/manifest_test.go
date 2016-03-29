package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	m, err := New("test", `{{ .test }}`)
	assert.Nil(t, err)

	out := m.Execute(map[string]string{"test": "hello!"})

	assert.Equal(t, "hello!", out)
}
