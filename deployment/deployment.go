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
)

var client coreclient.CoreClient
var deserializer runtime.Decoder

var groupVersionKind = unversioned.GroupVersionKind{
	Group:   api.GroupName,
	Version: runtime.APIVersionInternal,
	Kind:    meta.AnyKind,
}

func init() {

	scheme := runtime.NewScheme()
	v1.AddToScheme(scheme)
	factory := serializer.NewCodecFactory(scheme)
	deserializer = factory.UniversalDeserializer()

	kcfg := &restclient.Config{
		Host:     "http://localhost:8080",
		Insecure: true,
	}
	client, err := coreclient.NewForConfig(kcfg)
	if err != nil {
		panic(err)
	}
}

// Step represents a deployment step
type Step interface {
	Deploy() error
}

// DefaultStep implements a deployment step
type DefaultStep struct {
	object runtime.Object
}

var _ Step = &DefaultStep{}

// NewDefaultStep creates a default step
func NewDefaultStep(manifest string) (*DefaultStep, error) {
	object, _, err := deserializer.Decode([]byte(manifest), &groupVersionKind, nil)
	if err != nil {
		return nil, err
	}
	s := &DefaultStep{
		object: object,
	}
	return s, nil
}

// Deploy executes the deployment of a step
func (s *DefaultStep) Deploy() error {
	oGVK := s.object.GetObjectKind().GroupVersionKind()

	if oGVK.Kind == "ReplicationController" {
		rc, err := client.ReplicationControllers("default").Create(object.(*v1.ReplicationController))
		if err != nil {
			return err
		}
	}
	return nil
}
