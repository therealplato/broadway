package instance

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/namely/broadway/store"
)

// NotFoundError instance not found error
type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("broadway/instance: %s was not found", string(e))
}

// ErrMalformedSaveData in case an instance cannot be marshall
var ErrMalformedSaveData = errors.New("broadway/instance: saved data for this instance is malformed")

// Path represents a path for an instance
type Path struct {
	RootPath   string
	PlaybookID string
	ID         string
}

func (p Path) String() string {
	return fmt.Sprintf("%s/instances/%s/%s", p.RootPath, p.PlaybookID, p.ID)
}

// PlaybookPath represents a path for a playbook
type PlaybookPath struct {
	RootPath   string
	PlaybookID string
}

func (pp PlaybookPath) String() string {
	return fmt.Sprintf("%s/instances/%s", pp.RootPath, pp.PlaybookID)
}

// Instance entity
type Instance struct {
	PlaybookID string            `json:"playbook_id" binding:"required"`
	ID         string            `json:"id"`
	Created    int64             `json:"created_time"`
	Vars       map[string]string `json:"vars"`
	Status     `json:"status"`
	Path
}

// Status for an instance
type Status string

const (
	// StatusNew represents a newly created instance
	StatusNew Status = ""
	// StatusDeploying represents an instance that has begun deployment
	StatusDeploying = "deploying"
	// StatusDeployed represents an instance that has been deployed successfully
	StatusDeployed = "deployed"
	// StatusDeleting represents an instance that has begun deltion
	StatusDeleting = "deleting"
	// StatusError represents an instance that broke
	StatusError = "error"
)

// JSON instance representation
func (i *Instance) JSON() (string, error) {
	encoded, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// FindByPath find an instance based on it's path
func FindByPath(store store.Store, path Path) (*Instance, error) {
	i := store.Value(path.String())
	if i == "" {
		return nil, NotFoundError(path.String())
	}
	instance, err := fromJSON(i)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// FindByPlaybookID find all the instances for an specified playbook path
func FindByPlaybookID(store store.Store, playbookPath PlaybookPath) ([]*Instance, error) {
	data := store.Values(playbookPath.String())
	instances := []*Instance{}
	for _, value := range data {
		instance, err := fromJSON(value)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// Save an instance into the Store
func Save(store store.Store, instance *Instance) error {
	encoded, err := toJSON(instance)
	if err != nil {
		return err
	}
	err = store.SetValue(instance.Path.String(), encoded)
	return err
}

// Delete an instance from the store
func Delete(store store.Store, path Path) error {
	return store.Delete(path.String())
}

func fromJSON(jsonData string) (*Instance, error) {
	var instance Instance
	err := json.Unmarshal([]byte(jsonData), &instance)
	if err != nil {
		return nil, ErrMalformedSaveData
	}
	return &instance, nil
}

func toJSON(instance *Instance) (string, error) {
	encoded, err := json.Marshal(instance)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
