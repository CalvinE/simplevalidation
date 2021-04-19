package validator

import "reflect"

type Validateable interface {
	GetValidationErrors() map[string]string
}

type Validator interface {
	Validate(value interface{}, fieldName string, fieldKind reflect.Kind) (bool, error)
	ReadOptionsFromTagItems([]string) error
}

type ValidatorFactory func() Validator
