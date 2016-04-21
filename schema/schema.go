//package schema facilitates the binding and validation of a Model (from package github.com/laahs/gopkg/apify/resource)
package schema

// ----- Schema -----
// ------------------

//schema struct defines all the constraints and behaviours of the fields of a Modeler.
//it is useful for data validation, compliance checks, fields restrictions etc.
type Schema struct {
	Fields Fields
}

//type Field defines the constraints and behaviour of a specific field of a Modeler.
type Field struct {
}

//type Fields is an array of Field structs, created for method attachement.
type Fields []Field
