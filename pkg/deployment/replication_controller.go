package deployment

import (
	"errors"
	"time"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
)

// deployRC deploys the RC after it calls deleteRC
func deployRC(s *ManifestStep) error {
	var o *v1.ReplicationController
	switch s.object.(type) {
	case *v1.ReplicationController:
		o = s.object.(*v1.ReplicationController)
	case *api.ReplicationController:
		rr := s.object.(*api.ReplicationController)
		if err := scheme.Convert(rr, o); err != nil {
			glog.Error("API object conversion failed.")
			return err
		}
	}

	if rc, err := client.ReplicationControllers(namespace).Get(o.ObjectMeta.Name); err == nil && rc != nil {
		if compareRCs(rc, o) {
			glog.Info("Existing RC is identical, skipping deployment")
			return nil
		}

		if err := deleteRC(namespace, o.ObjectMeta.Name); err != nil {
			glog.Error(err)
		}
	}

	glog.Info("Creating new replication controller: ", o.ObjectMeta.Name)
	if _, err := client.ReplicationControllers(namespace).Create(o); err != nil {
		glog.Error("Create or Update failed: ", err)
		return err
	}

	return nil
}

// deleteRC scales down an RC and then deletes it
func deleteRC(namespace, metaName string) error {
	// SCALE RC DOWN TO 0
	rc, err := client.ReplicationControllers(namespace).Get(metaName)
	if err != nil {
		return err
	}
	// The i variable needs to be declared as a int32 for the Replicas type
	var i int32
	rc.Spec.Replicas = &i // Replicas type is *int32 ... so this is *int32(0)
	client.ReplicationControllers(namespace).Update(rc)
	time.Sleep(10 * time.Second) // Wait for Kubernetes to delete pods
	rc, err = client.ReplicationControllers(namespace).Get(metaName)
	if err != nil {
		return err
	}
	if rc.Status.Replicas != 0 {
		return errors.New("deployment: RC did not scale successfully")
	}
	// WATCH REPLICATION CONTROLLER
	// selector := fields.Set{"metadata.name": meta.Name}.AsSelector()
	// lo := api.ListOptions{Watch: true, FieldSelector: selector}
	// watcher, err := client.ReplicationControllers(namespace).Watch(lo)
	// defer watcher.Stop()
	//
	// var rc1 *v1.ReplicationController
	// var attempt int
	// for {
	// 	var ok bool
	// 	event := <-watcher.ResultChan()
	// 	rc1, ok = event.Object.(*v1.ReplicationController)
	// 	if !ok {
	// 		rcv := v1.ReplicationController{}
	// 		apirc := event.Object.(*api.ReplicationController)
	// 		if err = api.Scheme.Convert(apirc, &rcv); err != nil {
	// 			glog.Errorln("API Object conversion failed: ", err)
	// 			return err
	// 		}
	// 		rc1 = &rcv
	// 	}
	//
	// 	if rc1.Status.Replicas == 0 {
	// 		break
	// 	}
	// 	if attempt > 20 {
	// 		return errors.New("deployment: RC deletion timed out")
	// 	}
	// 	time.Sleep(200 * time.Millisecond)
	// 	attempt++
	// }
	// DELETE RC
	return client.ReplicationControllers(namespace).Delete(metaName, nil)
}
