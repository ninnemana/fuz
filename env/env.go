// env provides a way to get, list, set, etc. all environment variables in the local system.

package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// Option defines the requirements for each environemnt configuration object.
type Option interface {
	Validate() error
}

// AlgoliaAppID implements Option for the Algolia Application Identifier.
type AlgoliaAppID string

// Validate implements the validation rules for an Algolia APP ID.
func (a AlgoliaAppID) Validate() error {
	fmt.Println(a)
	return nil
	// return errors.New("not implemented")
}

// Configuration holds all environment variables that are available through the system.
// It uses struct tagging to specify the naming, requirements, etc. All types must also implement
// the Option interface to adhere to validation rules.
type Configuration struct {
	AlgoliaAppID *AlgoliaAppID `json:"algoliaAppID" osname:"ALGOLIA_APP_ID"`
	sync.Mutex
}

// New returns the default configuration options for the current environment.
func New(path string) (*Configuration, error) {
	var c Configuration
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}

		path = fmt.Sprintf("%s/env_test.json", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	val := reflect.ValueOf(c)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		field.IsNil
		switch field.Type() {
		case nil:
			fmt.Println(field.Type(), "nil")
		default:
			fmt.Println(reflect.TypeOf(field))
		}
	}

	return &c, nil
}

// Set updates a configration variable
func (c *Configuration) Set(o Option) error {
	err := o.Validate()
	if err != nil {
		return err
	}
	c.Lock()
	defer c.Unlock()

	fmt.Println(reflect.TypeOf(o).Name())
	return errors.New("not implemented")
}

// Clear empties all values
func (c *Configuration) Clear() error {
	c.Lock()
	defer c.Unlock()

	c = &Configuration{}

	return nil
}

// Validate goes over each property of the current configuration
// and validates it's value.
func (c *Configuration) Validate() error {
	val := reflect.ValueOf(c)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fmt.Println(reflect.TypeOf(field))
	}

	return errors.New("not implemented")
}
