package deployment

import (
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
)

// ManifestStep implements a deployment step
type ManifestStep struct {
	object runtime.Object
}

var _ Step = &ManifestStep{}

// NewManifestStep creates a default step
func NewManifestStep(manifest string) (Step, error) {
	object, _, err := deserializer.Decode([]byte(manifest), &groupVersionKind, nil)
	if err != nil {
		return nil, err
	}
	s := &ManifestStep{
		object: object,
	}
	return s, nil
}

// Deploy executes the deployment of a step
func (s *ManifestStep) Deploy() error {
	oGVK := s.object.GetObjectKind().GroupVersionKind()
	if oGVK.Kind == "ReplicationController" {
		rc := s.object.(*v1.ReplicationController)
		_, err := client.ReplicationControllers("default").Create(rc)
		if err != nil {
			return err
		}
	}
	return nil
}
