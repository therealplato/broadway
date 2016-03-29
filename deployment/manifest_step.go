package deployment

import (
	"errors"
	"log"

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
		o := s.object.(*v1.ReplicationController)

		_, err := client.ReplicationControllers(namespace).Get(o.ObjectMeta.Name)

		if err != nil {
			log.Println("Creating new replication controller: ", o.ObjectMeta.Name)
			_, err = client.ReplicationControllers(namespace).Create(o)
		} else {
			log.Println("Updating replication controller: ", o.ObjectMeta.Name)
			_, err = client.ReplicationControllers(namespace).Update(o)
		}
		if err != nil {
			log.Println("Create or Update failed: ", err)
			return err
		}
	case "Pod":
		o := s.object.(*v1.Pod)
		_, err := client.Pods(namespace).Get(o.ObjectMeta.Name)

		if err != nil {
			log.Println("Creating new pod: ", o.ObjectMeta.Name)
			_, err = client.Pods(namespace).Create(o)
		} else {
			log.Println("Updating pod", o.ObjectMeta.Name)
			_, err = client.Pods(namespace).Update(o)
		}
		if err != nil {
			log.Println("Create or Update failed: ", err)
			return err
		}
	case "Service":
		o := s.object.(*v1.Service)
		_, err := client.Services(namespace).Get(o.ObjectMeta.Name)

		if err != nil {
			log.Println("Creating new service: ", o.ObjectMeta.Name)
			_, err = client.Services(namespace).Create(o)
		} else {
			log.Println("Updating service", o.ObjectMeta.Name)
			_, err = client.Services(namespace).Update(o)
		}
		if err != nil {
			log.Println("Create or Update failed: ", err)
			return err
		}
	default:
		return errors.New("Kubernetes resource is not recognized: " + oGVK.Kind)
	}
	return nil
}
