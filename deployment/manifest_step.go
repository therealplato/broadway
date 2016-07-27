package deployment

import (
	"errors"
	"reflect"
	"time"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
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

func check(name string, r bool) bool {
	if r == false {
		glog.Info("Found difference: " + name)
	}
	return r
}

func compareI(name string, a, b interface{}) bool {
	return check(name, reflect.DeepEqual(a, b))
}

func compareS(name, a, b string) bool {
	return check(name, a == b)
}

func compareContainers(a, b v1.Container) bool {
	return a.Name == b.Name && a.Image == b.Image &&
		compareI("container command", a.Command, b.Command) &&
		compareI("container args", a.Args, b.Args) &&
		compareS("container workdir", a.WorkingDir, b.WorkingDir) &&
		compareI("container ports", a.Ports, b.Ports) &&
		compareI("container env", a.Env, b.Env) &&
		compareI("container resources", a.Resources, b.Resources) &&
		check("container pull policy", a.ImagePullPolicy == b.ImagePullPolicy)
}

func compareContainerLists(a, b []v1.Container) bool {
	aMap := map[string]v1.Container{}
	bMap := map[string]v1.Container{}

	for _, c := range a {
		aMap[c.Name] = c
	}
	for _, c := range b {
		bMap[c.Name] = c
	}

	if !check("conainer list len", len(aMap) == len(bMap)) {
		return false
	}

	for name, ac := range aMap {
		bc, ok := bMap[name]
		if !ok {
			return false
		}
		if !compareContainers(ac, bc) {
			return false
		}
	}
	return true
}

func comparePodSpecs(a, b v1.PodSpec) bool {
	if !check("pod spec containers len", len(a.Containers) == len(b.Containers)) {
		return false
	}
	if !compareContainerLists(a.Containers, b.Containers) {
		return false
	}
	return compareI("pod spec image pull secrets", a.ImagePullSecrets, b.ImagePullSecrets)
}

// compareRCs compares RC a to RC b. If a == b the return is true
func compareRCs(a, b *v1.ReplicationController) bool {
	if a.ObjectMeta.Name == "" {
		return false
	}
	return compareS("rc object meta name", a.ObjectMeta.Name, b.ObjectMeta.Name) &&
		compareI("rc object meta labels", a.ObjectMeta.Labels, b.ObjectMeta.Labels) &&
		compareI("rc spec replicas", a.Spec.Replicas, b.Spec.Replicas) &&
		(len(a.Spec.Selector) == 0 || len(b.Spec.Selector) == 0 || compareI("rc spec selector", a.Spec.Selector, b.Spec.Selector)) &&
		compareI("rc spec template object meta labels", a.Spec.Template.ObjectMeta.Labels, b.Spec.Template.ObjectMeta.Labels) &&
		compareS("rc spec template object meta name", a.Spec.Template.ObjectMeta.Name, b.Spec.Template.ObjectMeta.Name) &&
		comparePodSpecs(a.Spec.Template.Spec, b.Spec.Template.Spec)
}

func comparePods(a, b *v1.Pod) bool {
	return compareS("pod object meta name", a.ObjectMeta.Name, b.ObjectMeta.Name) &&
		compareI("pod object meta labels", a.ObjectMeta.Labels, b.ObjectMeta.Labels) &&
		comparePodSpecs(a.Spec, b.Spec)
}

// Deploy executes the deployment of a step
func (s *ManifestStep) Deploy() error {
	oGVK := s.object.GetObjectKind().GroupVersionKind()
	switch oGVK.Kind {
	case "ReplicationController":
		return deployRC(s)
	case "Pod":
		o := s.object.(*v1.Pod)
		pod, err := client.Pods(namespace).Get(o.ObjectMeta.Name)

		if err == nil && pod != nil {
			if comparePods(pod, o) {
				glog.Info("Existing Pod is identical, skipping deployment")
				return nil
			}
			glog.Info("Deleting old pod", o.ObjectMeta.Name)
			err = client.Pods(namespace).Delete(o.ObjectMeta.Name, nil)

			for k := 1; err == nil && k < 20; k++ {
				time.Sleep(200 * time.Millisecond) // Wait for Kubernetes to delete the resource
				_, err = client.Pods(namespace).Get(o.ObjectMeta.Name)
			}
			if err != nil {
				glog.Error("delete old pods: ", err)
			}
		}

		glog.Info("Creating new pod: ", o.ObjectMeta.Name)
		_, err = client.Pods(namespace).Create(o)
		if err != nil {
			glog.Info("Create or Update failed: ", err)
			return err
		}
	case "Service":
		o := s.object.(*v1.Service)
		service, err := client.Services(namespace).Get(o.ObjectMeta.Name)

		if err != nil {
			glog.Info("Creating new service: ", o.ObjectMeta.Name)
			_, err = client.Services(namespace).Create(o)
		} else {
			glog.Info("Updating service", o.ObjectMeta.Name)
			o.ObjectMeta.ResourceVersion = service.ObjectMeta.ResourceVersion
			o.Spec.ClusterIP = service.Spec.ClusterIP
			_, err = client.Services(namespace).Update(o)
		}
		if err != nil {
			glog.Info("Create or Update failed: ", err)
			return err
		}
	default:
		return errors.New("Kubernetes resource is not recognized: " + oGVK.Kind)
	}
	return nil
}

// Destroy deletes kubernetes resource
func (s *ManifestStep) Destroy() error {
	var err error
	oGVK := s.object.GetObjectKind().GroupVersionKind()
	meta, err := api.ObjectMetaFor(s.object)
	if err != nil {
		return err
	}
	switch oGVK.Kind {
	case "ReplicationController":
		err = deleteRC(namespace, meta.Name)
	case "Service":
		err = client.Services(namespace).Delete(meta.Name, nil)
	case "Pod":
		err = client.Pods(namespace).Delete(meta.Name, nil)
	}

	return err
}
