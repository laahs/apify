//package api defines the api and methods for rest requests handling.
package resource

//Package apify helps get standard structure for api definition, playing with models via http requests
//to a rest api. From

import (
	"errors"
	"github.com/laahs/gopkg/apify/schema"
	"regexp"
	"strings"
	"sync"
	//"github.com/rs/xlog"
	//"golang.org/x/net/context"
)

//type Name is a string empowered with auto validation. useful to ensure that all names satisfy same requirements. It is used in
//github.com/laahs/gopkg/apify/api to validate names in targets definition.
type Name string

func (n *Name) IsValid() (bool, error) {
	n = strings.TrimSpace(strings.ToLower(n))
	matched, err := regexp.MatchString("^[a-z]{4,15}$", n)
	if err == nil && matched {
		return nil
	} else {
		return errors.New("parent resource name does not match requirements and has not been set: must be ^[a-z]{4-15}$")
	}
}

// ----- Resource -----
// -----------------

//type Model is a structure containing all the information about the underlying item to be store (Item interface),
//the store (Storer interface) where to store the data, other models it is related to and that it depends on, the http
//method that are allowed on the attributed endpoint, and the schema for automatic data validation and restrictions.
//the description is used by the api documentation generator.
type Resource struct {
	name           string
	description    string
	schema         *Schema
	model          *Modeler
	parentresource string //list of models names that this model is related to and depends on. must contain theses dependencies names as fields. useful for one-shot queries like "users/:id/posts/:id"
	store          *Storer
	methods        Methods
	path           string //endpoint to send request for this resource. automatically generated depending on parent (chaining parents paths from root, to undirect parents then direct parent). format '/undirectparent/:id/directparent/:id/name' .
}

// ----- parentresource -------

//method SetParent of Resource sets parent model to the current model, to be able to specify dependency and define paths chaining models.
func (m *Resource) SetParent(name string) error {
	if len(m.name) > 0 {
		return errors.New("name already set")
	}
	n := Name(name)
	ok, err := n.IsValid()
	if err == nil && ok {
		m.parentresource = n
		return nil
	} else {
		return err
	}
}

//method Name of Resource returns the name attached to the model, it is used for api doc generation and automatic routes definition. name is an unexported field.
func (m *Resource) Parent() string {
	return m.parentresource
}

// ----- path -------

//method SetName of Resource sets the name attached to the model, to be able to call it by name from the Models object. name is an unexported field.
func (m *Resource) SetPath(name string) error {
	if len(m.name) > 0 {
		return errors.New("name already set")
	}
	name = strings.TrimSpace(strings.ToLower(s))
	matched, err := regexp.MatchString("^[a-z]{4,15}$", name)
	if err == nil && matched {
		m.name = name
		return nil
	} else {
		return errors.New("name does not match requirements and has not been set: must be ^[a-z]{4-15}$")
	}
}

//method Name of Resource returns the name attached to the model, it is used for api doc generation and automatic routes definition. name is an unexported field.
func (m *Resource) Path() string {
	return m.path
}

// ----- name -------

//method SetName of Resource sets the name attached to the model, to be able to call it by name from the Models object. name is an unexported field.
func (m *Resource) SetName(name string) error {
	if len(m.name) > 0 {
		return errors.New("name already set")
	}
	n := Name(name)
	ok, err := n.IsValid()
	if err == nil && ok {
		m.name = n
		return nil
	} else {
		return err
	}
}

//method Name of Resource returns the name attached to the model, it is used for api doc generation and automatic routes definition. name is an unexported field.
func (m *Resource) Name() string {
	return m.name
}

// ----- description -----

//method SetDescription of Resource sets the description attached to the model. description is an unexported field.
func (m *Resource) SetDescription(descr string) error {
	if len(m.description) > 0 {
		return errors.New("description already set")
	}
	if len(descr) > 15 && len(descr) < 200 {
		return nil
	} else {
		return errors.New("description does not match requirements and has not been set: must contain 15<caracters<200")
	}
}

//method Description of Resource returns the description attached to the model, it is used by the api documentation generator. description is an unexported field.
func (m *Resource) Description() string {
	return m.description
}

// ----- Resource validation -----

//method IsOK of Resource returns processes all the checks to validate Resource definition.
func (m *Resource) IsOK() bool {
	switch {
	case len(m.name) == 0:
		return false
	case len(m.description) == 0:
		return false
	case m.schema == nil:
		return false
	case m.model == nil:
		return false
	case m.store == nil:
		return false
	case len(m.methods) == 0:
		return false
	default:
		return true
	}
}

// ----- method -----
// ------------------

//type Method is a maps http method names to restrictions. Useful for basic and automatic access
//restrictions (used by auth middlewares for instance).
type Method struct {
	meth         string
	restrictions HttpRestrictions
	//mx           sync.Mutex
}

//method Method of Method
func (m *Method) Method() {

}

//method Validate of Method validates if the all the keys match

// ----- methods -----
// ------------------

//type Methods is an array of method maps.
type Methods []Method

//method
// ----- Resources -----
// ------------------

//type Resources is a list of Resource, store in the api definition. Created for easier method call.
type Resources struct {
	list []string
}

//method Names of Resources returns the list of all the names of the resources.
func (rs *Resources) Names() []string {
	names := []string{}
	for _, r := range rs {
		names = append(names, r.list.Name())
	}
	return names
}

// ----- Modeler -----
// ------------------

//interface Model is a generic type representing the data that should be manipulated and store.
type Modeler interface {
	GetStore() Store
}
