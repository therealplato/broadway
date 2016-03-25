package deployment

import (
	"errors"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
)

// ManifestStep implements a deployment step
type ManifestStep struct {
	object runtime.Object
}

var _ Step = &ManifestStep{}

// NewManifestStep creates a default step
func NewManifestStep(object runtime.Object) Step {
	return &ManifestStep{
		object: object,
	}
}

// Deploy executes the deployment of a step
func (s *ManifestStep) Deploy() error {
	oGVK := s.object.GetObjectKind().GroupVersionKind()
	switch oGVK.Kind {
	case "ReplicationController":
		rc := s.object.(*v1.ReplicationController)
		_, err := client.ReplicationControllers(namespace).Create(rc)
		if err != nil {
			return err
		}
	case "Pod":
		pod := s.object.(*v1.Pod)
		_, err := client.Pods(namespace).Create(pod)
		if err != nil {
			return err
		}
	case "Service":
		service := s.object.(*v1.Service)
		_, err := client.Services(namespace).Create(service)
		if err != nil {
			return err
		}
	default:
		return errors.New("Manifest is not recognized")
	}
	return nil
}
