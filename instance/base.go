package instance

import (
	"encoding/json"
	"errors"

	"github.com/namely/broadway/store"
)

type baseInstance struct {
	attributes *InstanceAttributes

	store store.Store
}

// New constructs an Instance from a store and attributes
func New(s store.Store, attrs *InstanceAttributes) Instance {
	instance := &baseInstance{
		attributes: attrs,
		store:      s,
	}

	return instance
}

// Get looks up an Instance by playbookId and instance id
func Get(playbookId, id string) (Instance, error) {
	attrs := &InstanceAttributes{
		PlaybookId: playbookId,
		Id:         id,
	}

	instance := &baseInstance{
		attributes: attrs,
		store:      store.New(),
	}

	value := instance.store.Value(instance.path())
	if value == "" {
		return nil, errors.New("Instance does not exist.")
	}
	err := json.Unmarshal([]byte(value), instance.attributes)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// List looks up all Instances stored under a given playbookId
func List(s store.Store, playbookId string) ([]Instance, error) {
	instances := []Instance{}
	instanceKeysValues := s.Values("/broadway/instances/" + playbookId)
	for _, v := range instanceKeysValues {
		var attrs InstanceAttributes
		err := json.Unmarshal([]byte(v), &attrs)
		if err != nil {
			return nil, err
		}
		i := &baseInstance{
			attributes: &attrs,
			store:      s,
		}
		instances = append(instances, i)
	}
	return instances, nil
}

// ID returns the instance id
func (instance *baseInstance) ID() string {
	return instance.Attributes().Id
}

// PlaybookID returns the instance's playbook id
func (instance *baseInstance) PlaybookID() string {
	return instance.Attributes().PlaybookId
}

// Status returns the instance status
func (instance *baseInstance) Status() InstanceStatus {
	return instance.Attributes().Status
}

// Attributes returns the instance attributes
func (instance *baseInstance) Attributes() *InstanceAttributes {
	return instance.attributes
}

// MarshalJSON implements the json.Marshaler interface. Only the instance's
// attributes are serialized.
func (instance *baseInstance) MarshalJSON() ([]byte, error) {
	o, err := instance.Attributes().JSON()
	if err != nil {
		return nil, err
	}
	return []byte(o), nil
}

func (instance *baseInstance) path() string {
	return "/broadway/instances/" + instance.PlaybookID() + "/" + instance.ID()
}

// Save sets or updates the stored instance, keyed on playbookId and instance id
func (instance *baseInstance) Save() (err error) {
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
func (instance *baseInstance) Destroy() error {
	return instance.store.Delete(instance.path())
}
