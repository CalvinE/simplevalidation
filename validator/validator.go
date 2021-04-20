package validator

import "reflect"

const (
	InvalidTypeErrorTemplate = "type: the value of %s is of type %T which is not valid"
)

// This interface represents a data validator, they are intended to be simple and easy to implement for custom types.
type Validator interface {
	// Validate takes a value, fieldName if from struct, and the fieldKind and returns true if valid and false if not.
	// If the value is not valid then an error should also be returned with info on why its invalid.
	Validate(value interface{}, fieldName string, fieldKind reflect.Kind) (bool, error)
	/*
		ReadOptionsFromTagItems takes in an array of tag arguments and reads them into the validator.
		This function is intended to be attached to a pointer of the validator,
		and the pointer has its options set from this method.
		An error should be returned if there is an issue populating the Validator options.

		For more information on the validation tag syntax please see the package documentation for the validtor package.

		An example of the implementation signature would be:

			func (v *myCustomValidator) ReadOptionsFromTagItems error
	*/
	ReadOptionsFromTagItems([]string) error
}

// ValidatorFactory is a type alias for a function that take no parameters and reutrns a Validator.
type ValidatorFactory func() Validator
