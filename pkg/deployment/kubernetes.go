package deployment

import (
	"github.com/golang/glog"
	"github.com/namely/broadway/pkg/cfg"

	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	coreclient "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3/typed/core/v1"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/runtime/serializer"

	// Install API
	_ "k8s.io/kubernetes/pkg/api/install"
)

var groupVersionKind = unversioned.GroupVersionKind{
	Group:   v1.GroupName,
	Version: runtime.APIVersionInternal,
	Kind:    meta.AnyKind,
}

var client coreclient.CoreInterface
var deserializer runtime.Decoder
var namespace string
var scheme *runtime.Scheme

// Step represents a deployment step
type Step interface {
	Deploy() error
	Destroy() error
}

// SetupKubernetes configures kubernetes with an injected configuration
func SetupKubernetes(cfg cfg.Type) {
	scheme = runtime.NewScheme()
	v1.AddToScheme(scheme)
	factory := serializer.NewCodecFactory(scheme)
	deserializer = factory.UniversalDeserializer()

	namespace = cfg.K8sNamespace
}

// KubernetesDeployment represents a deployment of an instance
type KubernetesDeployment struct {
	Playbook  *Playbook
	Variables map[string]string
	Manifests map[string]*Manifest
}

// NewKubernetesDeployment creates a new kuberentes deployment
func NewKubernetesDeployment(config *restclient.Config, playbook *Playbook, variables map[string]string, manifests map[string]*Manifest) (*KubernetesDeployment, error) {
	var err error
	vars := map[string]string{}
	client, err = coreclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	for _, v := range playbook.Vars {
		var ok bool
		vars[v], ok = variables[v]
		if !ok {
			vars[v] = ""
		}
	}
	return &KubernetesDeployment{
		Playbook:  playbook,
		Variables: vars,
		Manifests: manifests,
	}, nil
}

// Deploy executes the deployment
func (d *KubernetesDeployment) Deploy() error {
	steps, err := d.steps()
	if err != nil {
		return err
	}

	for i, step := range steps {
		err := step.Deploy()
		if err != nil {
			glog.Warning("%d. step failed: %s", i, err.Error())
			return err
		}
	}

	return nil
}

// Destroy deletes Kubernetes resourses
func (d *KubernetesDeployment) Destroy() error {
	steps, err := d.steps()
	if err != nil {
		return err
	}

	for i, step := range steps {
		glog.Infof("%d. Destroying Resources.", i)
		err := step.Destroy()
		if err != nil {
			glog.Warning("%d. step failed: %s", i, err.Error())
			return err
		}
	}
	return nil
}

func (d *KubernetesDeployment) steps() ([]Step, error) {
	var steps = []Step{}
	for _, name := range d.Playbook.Manifests {
		m := d.Manifests[name]
		rendered := m.Execute(d.Variables)
		object, err := deserialize(rendered)
		if err != nil {
			glog.Warningf("Failed to parse manifest %s", name)
			return steps, err
		}
		steps = append(steps, NewManifestStep(object))
	}
	return steps, nil
}

func deserialize(manifest string) (runtime.Object, error) {
	object, _, err := deserializer.Decode([]byte(manifest), &groupVersionKind, nil)
	if err != nil {
		return nil, err
	}
	return object, nil
}
