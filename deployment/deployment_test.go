package deployment

import (
	"testing"

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
	cases := []struct {
		Tasks    []playbook.Task
		Expected int
	}{
		{
			Tasks: []playbook.Task{
				{
					Name: "First step",
					Manifests: []string{
						"test",
					},
				},
			},
			Expected: 1,
		}, {
			Tasks: []playbook.Task{
				{
					Name: "First step",
					Manifests: []string{
						"test",
						"test2",
						"test2",
					},
				},
			},
			Expected: 3,
		}, {
			Tasks: []playbook.Task{
				{
					Name:        "First step",
					PodManifest: "test",
				},
			},
			Expected: 1,
		},
	}

	vars := map[string]string{
		"test": "ok",
	}
	m, _ := manifest.New("test", mtemplate)
	manifests := map[string]*manifest.Manifest{
		"test":  m,
		"test2": m,
		"test3": m,
	}

	for _, c := range cases {
		// Reset client
		client.(*fake.FakeCore).Fake.ClearActions()

		p := playbook.Playbook{
			ID:    "test",
			Name:  "Test deployment",
			Meta:  playbook.Meta{},
			Vars:  []string{"test"},
			Tasks: c.Tasks,
		}

		d := &Deployment{
			Playbook:  p,
			Variables: vars,
			Manifests: manifests,
		}

		err := d.Deploy()
		assert.Nil(t, err)
		f := client.(*fake.FakeCore).Fake
		assert.Equal(t, c.Expected, len(f.Actions()))
		for _, action := range f.Actions() {
			assert.IsType(t, core.CreateActionImpl{}, action)
		}
	}
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
