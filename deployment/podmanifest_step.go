package deployment

import "k8s.io/kubernetes/pkg/runtime"

// PodmanifestStep implements a deployment step with pod manifest
type PodmanifestStep struct {
	object runtime.Object
}

var _ Step = &ManifestStep{}

// NewPodManifestStep creates a podmanifest step and returns a Step
func NewPodmanifestStep(podmanifest string) (Step, error) {
	object, _, err := deserializer.Decode([]byte(podmanifest), &groupVersionKind, nil)
	if err != nil {
		return nil, err
	}
	s := &ManifestStep{
		object: object,
	}
	return s, nil
}

// Deploy executes the deployment step
func (s *PodmanifestStep) Deploy() error {
	return nil
}
