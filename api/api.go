//package api defines automatically packs handlers, resources and schemas together to generate
//a rest api that handles http requests to a set of defined resources
package api

import (
	"errors"
	"github.com/laahs/gopkg/apify/resource"
	"strings"
)

//type Api is the most basic block of this package. API definition is attached to this object.
type Api struct {
	t targets //targets are the resources handled by the api.
	v string  //api version
}

//method AddTarget of Api adds a target (a)
func (a *Api) AddTarget(name string, rsrc resource.Resource) error {
	name=
	if a.t[]
}

//type targets is a mapping resource names and resource definitions
type targets map[string]resource.Resource
