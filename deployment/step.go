package deployment

import (
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/namely/broadway/playbook"
)

// Step represents a deployment step
type Step interface {
	Task() playbook.Task
	Deploy() error
}

// DefaultStep implements a deployment step
type DefaultStep struct {
	task   playbook.Task
	object runtime.Object
}

var _ Step = &DefaultStep{}

// NewDefaultStep creates a default step
func NewDefaultStep(task playbook.Task, manifest string) (*DefaultStep, error) {
	object, _, err := deserializer.Decode([]byte(manifest), &groupVersionKind, nil)
	if err != nil {
		return nil, err
	}
	s := &DefaultStep{
		object: object,
		task:   task,
	}
	return s, nil
}

// Deploy executes the deployment of a step
func (s *DefaultStep) Deploy() error {
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

// Task returns the step task
func (s *DefaultStep) Task() playbook.Task {
	return s.task
}
