package deployment

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/client/restclient"
	coreclient "k8s.io/kubernetes/pkg/client/typed/generated/core/v1"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/runtime/serializer"

	// Install API
	_ "k8s.io/kubernetes/pkg/api/install"

	"github.com/namely/broadway/manifest"
	"github.com/namely/broadway/playbook"
)

var groupVersionKind = unversioned.GroupVersionKind{
	Group:   api.GroupName,
	Version: runtime.APIVersionInternal,
	Kind:    meta.AnyKind,
}

var client coreclient.CoreInterface
var deserializer runtime.Decoder

func init() {

	scheme := runtime.NewScheme()
	v1.AddToScheme(scheme)
	factory := serializer.NewCodecFactory(scheme)
	deserializer = factory.UniversalDeserializer()

	kcfg := &restclient.Config{
		Host:     "http://localhost:8080",
		Insecure: true,
	}
	var err error
	client, err = coreclient.NewForConfig(kcfg)
	if err != nil {
		panic(err)
	}
}

// Deployer declares something that can Deploy Deployments
type Deployer interface {
	Deploy(playbook.Playbook, map[string]string) error
}

// Deployment represents a deployment of an instance
type Deployment struct {
	Playbook  playbook.Playbook
	Variables map[string]string
	Manifests map[string]*manifest.Manifest
}

// Deploy executes the deployment
func (d *Deployment) Deploy() error {
	for _, task := range d.Playbook.Tasks {
		if len(task.Manifests) > 0 {
			for _, name := range task.Manifests {
				m := d.Manifests[name]
				rendered := m.Execute(d.Variables)
				step, err := NewDefaultStep(task, rendered)
				if err != nil {
					return err
				}
				err = step.Deploy()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
