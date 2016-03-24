package instance

import (
	"encoding/json"

	"github.com/namely/broadway/store"
)

type defaultInstance struct {
	attributes *Attributes

	store store.Store
}

// New constructs an Instance from a store and attributes
func New(s store.Store, attrs *Attributes) Instance {
	instance := &defaultInstance{
		attributes: attrs,
		store:      s,
	}

	return instance
}

// List looks up all Instances stored under a given playbookID
func List(s store.Store, playbookID string) ([]Instance, error) {
	instances := []Instance{}
	instanceKeysValues := s.Values("/broadway/instances/" + playbookID)
	for _, v := range instanceKeysValues {
		var attrs Attributes
		err := json.Unmarshal([]byte(v), &attrs)
		if err != nil {
			return nil, err
		}
		i := &defaultInstance{
			attributes: &attrs,
			store:      s,
		}
		instances = append(instances, i)
	}
	return instances, nil
}

// ID returns the instance id
func (instance *defaultInstance) ID() string {
	return instance.Attributes().ID
}

// PlaybookID returns the instance's playbook id
func (instance *defaultInstance) PlaybookID() string {
	return instance.Attributes().PlaybookID
}

// Status returns the instance status
func (instance *defaultInstance) Status() Status {
	return instance.Attributes().Status
}

// Attributes returns the instance attributes
func (instance *defaultInstance) Attributes() *Attributes {
	return instance.attributes
}

// MarshalJSON implements the json.Marshaler interface. Only the instance's
// attributes are serialized.
func (instance *defaultInstance) MarshalJSON() ([]byte, error) {
	o, err := instance.Attributes().JSON()
	if err != nil {
		return nil, err
	}
	return []byte(o), nil
}

func (instance *defaultInstance) path() string {
	return "/broadway/instances/" + instance.PlaybookID() + "/" + instance.ID()
}

// Save sets or updates the stored instance, keyed on playbookID and instance id
func (instance *defaultInstance) Save() (err error) {
	encoded, err := instance.Attributes().JSON()
	if err != nil {
		return err
	}
	err = instance.store.SetValue(instance.path(), encoded)
	if err != nil {
		return err
	}
	return nil
}

// Destroy removes the stored instance
func (instance *defaultInstance) Destroy() error {
	return instance.store.Delete(instance.path())
}
