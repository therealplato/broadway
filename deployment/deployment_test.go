package deployment

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/client/testing/core"
	"k8s.io/kubernetes/pkg/client/typed/generated/core/v1/fake"

	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
)

func init() {
	client = &fake.FakeCore{&core.Fake{}}
}

func TestDeploy(t *testing.T) {
	p := playbook.Playbook{
		ID:   "test",
		Name: "Test deployment",
		Meta: playbook.Meta{},
		Vars: []string{"test"},
		Tasks: []playbook.Task{
			{
				Name: "First step",
				Manifests: []string{
					"test",
				},
			},
		},
	}

	v := map[string]string{
		"test": "ok",
	}

	//m, _ := manifest.New("test", mtemplate)

	id := "test"
	tp, err := template.New(id).Parse(mtemplate)
	assert.Nil(t, err)
	m := &manifest.Manifest{ID: id, Template: tp}

	ms := map[string]*manifest.Manifest{
		"test": m,
	}

	d := &Deployment{
		Playbook:  p,
		Variables: v,
		Manifests: ms,
	}

	err = d.Deploy()
	assert.Nil(t, err)
	f := client.(*fake.FakeCore).Fake
	assert.Equal(t, 1, len(f.Actions()))
}

var mtemplate = `apiVersion: v1
kind: ReplicationController
metadata:
  name: test
spec:
  replicas: 1
  selector:
    name: redis
  template:
    metadata:
      labels:
        name: redis
    spec:
      containers:
      - name: redis
        image: kubernetes/redis:v1
`
