package deployment

import "k8s.io/kubernetes/pkg/runtime"

// PodmanifestStep implements a deployment step with pod manifest
type PodmanifestStep struct {
	object runtime.Object
}

var _ Step = &ManifestStep{}

// NewPodmanifestStep creates a podmanifest step and returns a Step
func NewPodmanifestStep(object runtime.Object) Step {
	return &ManifestStep{
		object: object,
	}
}

// Deploy executes the deployment step
func (s *PodmanifestStep) Deploy() error {
	return nil
}
