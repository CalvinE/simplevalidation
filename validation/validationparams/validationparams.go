package validationparams

import "github.com/calvine/simplevalidation/validator"

type ValidationParams struct {
	ArrayDepth     uint8
	FieldValidator validator.Validator
	Name           string
	Required       bool
	StructDepth    uint8
	Value          interface{}
}

func New() ValidationParams {
	return ValidationParams{
		ArrayDepth:  0,
		Name:        "",
		Required:    false,
		StructDepth: 0,
	}
}
