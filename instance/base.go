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

func New(s store.Store, attrs *InstanceAttributes) Instance {
	instance := &baseInstance{
		attributes: attrs,
		store:      s,
	}

	return instance
}

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

func (instance *baseInstance) ID() string {
	return instance.Attributes().Id
}

func (instance *baseInstance) PlaybookID() string {
	return instance.Attributes().PlaybookId
}

func (instance *baseInstance) Status() InstanceStatus {
	return instance.Attributes().Status
}

func (instance *baseInstance) Attributes() *InstanceAttributes {
	return instance.attributes
}

func (instance *baseInstance) path() string {
	return "/broadway/instances/" + instance.PlaybookID() + "/" + instance.ID()
}

func (instance *baseInstance) Save() (err error) {
	encoded, err := instance.Attributes().JSON()
	if err != nil {
		return err
	}
	instance.store.SetValue(instance.path(), encoded)
	return err
}

func (instance *baseInstance) Destroy() error {
	return instance.store.Delete(instance.path())
}
