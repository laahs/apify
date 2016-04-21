package resource

import (
	"net/http"
	"strings"
	"errors"
	"fmt"
)

// Conf defines the configuration for a given resource
type Conf struct {
	// AllowedModes is the list of Mode allowed for the resource.
	AllowedModes []Mode
	// DefaultPageSize defines a default number of items per page. By defaut,
	// no default page size is set resulting in no pagination if no `limit` parameter
	// is provided.
	PaginationDefaultLimit int
}

// Mode defines CRUDL modes to be used with Conf.AllowedModes.
type ModeCode int

const (
	// Create mode represents the POST method on a collection URL or the PUT method
	// on a _non-existing_ item URL.
	Create ModeCode = iota
	// Read mode represents the GET method on an item URL
	Read
	// Update mode represents the PATCH on an item URL.
	Update
	// Replace mode represents the PUT methods on an existing item URL.
	Replace
	// Delete mode represents the DELETE method on an item URL.
	Delete
	// Clear mode represents the DELETE method on a collection URL
	Clear
	// List mode represents the GET method on a collection URL.
	List
	// Search mode represents the GET method on a collection URL with query parameters (activating search engine).
	Search
)

var (
	// ReadWrite is a shortcut for all modes
	ReadWrite = []ModeCode{Create, Read, Update, Replace, Delete, List, Clear}
	// ReadOnly is a shortcut for Read and List modes
	ReadOnly = []ModeCode{Read, List}
	// WriteOnly is a shortcut for Create, Update, Delete modes
	WriteOnly = []ModeCode{Create, Update, Replace, Delete, Clear}

	// DefaultConf defines a configuration with some sensible default parameters.
	// Mode is read/write and default pagination limit is set to 20 items.
	DefaultConf = Conf{
		AllowedModes:           ReadWrite,
		PaginationDefaultLimit: 20,
	}
)

// IsModeAllowed returns true if the provided mode is allowed in the configuration
func (c Conf) IsModeAllowed(mode ModeCode) bool {
	for _, m := range c.AllowedModes {
		if m == mode {
			return true
		}
	}
	return false
}

//httpModeMatcher is an interface that is able to transcript the modecode to http request type and vice versa.
type HttpModeMatcher interface {
	ModeToHttpMethod(ModeCode) (Hmeth,error)
	HttpMethodToMode(http.Request.Method, withId bool,isSearch bool) (ModeCode,error)
}

//type Matcher is an empty struct implementing the HttpModeMatcher interface
type Matcher struct{}
//method ModeToHttpMethod analyses the ModeCode value and return the matching Http.Request.Method (if any), a boolean 
//specifying if action should be sent at root of resource (false) or with a resource id (true), a boolean specifying if
//query parameters are handled/required and an error if the ModeCode does not match any http.Request.Method .
func (m *Matcher) ModeToHttpMethod(code ModeCode) (Hmeth,error) {
	switch code{
	case Create:
		return Hmeth{"POST",false,false},nil
	case Read:
		return Hmeth{"GET",true,false},nil
	case Update:
		return Hmeth{"PATCH",true,false},nil
	case Replace:
		return Hmeth{"PUT",true,false},nil
	case Delete:
		return Hmeth{"DELETE",true,false},nil
	case Clear:
		return Hmeth{"DELETE",false,false},nil
	case List:
		return Hmeth{"GET",false,false},nil
	case Search:
		return Hmeth{"GET",false,true},nil
	default:
		return nil,errors.New("method not allowed")
	}
}
//method HttpMethodToMode decodes the set of arguments {method, withId, isSearch} and return the corresponding
//ModeCode.
func (m *Matcher) HttpMethodToMode(method http.Request.Method, withId bool,isSearch bool) (ModeCode,error) {
	switch method{
	case "GET":
		switch {
		case withId && !isSearch:
			return Read,nil
		case !withId && isSearch:
			return Search,nil
		case !withId && !isSearch:
			return List,nil
		default:
			msg:=fmt.Sprintf("can not handle %s request with pair 'withId'=%v and 'isSearch'=%v", method,withId,isSearch)
			return -1,errors.New(msg)
		}
	case "POST":
		switch {
		case !withId && !isSearch:
			return Create,nil
		default:
			msg:=fmt.Sprintf("can not handle %s request with pair 'withId'=%v and 'isSearch'=%v", method,withId,isSearch)
			return -1,errors.New(msg)
		}
	case "PUT":
		switch {
		case withId && !isSearch:
			return Replace,nil
		default:
			msg:=fmt.Sprintf("can not handle %s request with pair 'withId'=%v and 'isSearch'=%v", method,withId,isSearch)
			return -1,errors.New(msg)
		}
	case "PATCH":
		switch {
		case withId && !isSearch:
			return Update,nil
		default:
			msg:=fmt.Sprintf("can not handle %s request with pair 'withId'=%v and 'isSearch'=%v", method,withId,isSearch)
			return -1,errors.New(msg)
		}
	case "DELETE":
		switch {
		case withId && !isSearch:
			return Delete,nil
		case !withId && !isSearch:
			return Clear,nil
		default:
			msg:=fmt.Sprintf("can not handle %s request with pair 'withId'=%v and 'isSearch'=%v", method,withId,isSearch)
			return -1,errors.New(msg)
		}
	default:
		msg:=fmt.Sprintf("can not handle %s request with pair 'withId'=%v and 'isSearch'=%v", method,withId,isSearch)
		return -1,errors.New(msg)
	}
}

//type Hmeth is a struct used for method definition and a set of additional info
type Hmeth struct{
	Method string //http.Request.Method
	IdReq bool // is the resource id required or should the request be sent to the root of the resource path
	ParamReq bool //true for search engine mode
}
/*
func MethodNeedsId(s string) (bool,error){
	s=strings.ToUpper(s)
	switch s{
	case "GET":
	case "PUT":
		case ""
	}
}
*/
