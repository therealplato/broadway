package deployment

import (
	"testing"

	"k8s.io/kubernetes/pkg/client/testing/core"
	"k8s.io/kubernetes/pkg/client/typed/generated/core/v1/fake"
)

func init() {
	client = &fake.FakeCore{&core.Fake{}}
}

func TestManifestStepDeploy(t *testing.T) {
	//f := client.(*fake.FakeCore).Fake
	//cases := []struct {
	//	Object   runtime.Object
	//	Expected string
	//}{
	//	{
	//		Object:   &v1.ReplicationController{},
	//		Expected: "replicationController",
	//	},
	//}

	//for _, c := range cases {
	//	// Reset client
	//	client.(*fake.FakeCore).Fake.ClearActions()
	//	step := NewManifestStep(c.Object)
	//	err := step.Deploy()
	//	assert.Nil(t, err)
	//	assert.Equal(t, 1, len(f.Actions()))
	//	assert.Equal(t, "create", f.Actions()[0].GetVerb)
	//	assert.Equal(t, c.Expected, f.Actions()[0].GetResource())
	//}
}
