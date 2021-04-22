package uuidvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/calvine/simplevalidation/validator"
	"github.com/google/uuid"
)

// specifically the google uuid package uuid struct.
type uuidValidator struct {
	// AllowEmptyUUID when true will consider default (empty) uuid is valid, if not it will cause a validation error.
	AllowEmptyUUID bool
	// AllowString allows strings to be parsed into uuids for validation, if false a type validation error will be returned.
	AllowString bool
	// Required really only checks if a string is passed in as an empty string.
	Required bool
}

const (
	uuidStringProvidedErrorTemplate = "type: The field %s is a string, but allowstring was not provided in the validation tag"
	uuidInvalidStringErrorTemplate  = "invalid: the field %s has the value %s which could not be parsed into a UUID"
	uuidNoEmptyUUIDErrorTemplate    = "no empty: the field %s has an empty uuid value, but AllowEmptyUUID is false"
	uuidNoValueErrorTemplate        = "required: The field %s had an empty value: %v"
)

func New() validator.Validator {
	validator := &uuidValidator{}
	return validator
}

func (uv *uuidValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	var value uuid.UUID
	var stringValue string
	var parseError error
	emptyUUID := uuid.UUID{}
	noValueProvided := false
	emptyUUIDProvided := false
	switch t := n.(type) {
	case string:
		if !uv.AllowString {
			errorMessage := fmt.Sprintf(uuidStringProvidedErrorTemplate, fieldName)
			return false, errors.New(errorMessage)
		}
		stringValue = t
		if stringValue == "" {
			noValueProvided = true
		} else {
			value, parseError = uuid.Parse(t)
			if value == emptyUUID && parseError == nil {
				emptyUUIDProvided = true
			}
		}
	default:
		var ok bool
		value, ok = t.(uuid.UUID)
		if !ok {
			errorMessage := fmt.Sprintf(validator.InvalidTypeErrorTemplate, fieldName, t)
			return false, errors.New(errorMessage)
		}
		if t == emptyUUID {
			emptyUUIDProvided = true
		}
	}

	if uv.Required && noValueProvided {
		errorMessage := fmt.Sprintf(uuidNoValueErrorTemplate, fieldName, value)
		return false, errors.New(errorMessage)
	}

	if !uv.AllowEmptyUUID && emptyUUIDProvided {
		errorMessage := fmt.Sprintf(uuidNoEmptyUUIDErrorTemplate, fieldName)
		return false, errors.New(errorMessage)
	}

	if parseError != nil {
		errorMessage := fmt.Sprintf(uuidInvalidStringErrorTemplate, fieldName, stringValue)
		return false, errors.New(errorMessage)
	}
	return true, nil

}

func (uv *uuidValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "required":
			uv.Required = true
		case "allowemptyuuid":
			uv.AllowEmptyUUID = true
		case "allowstring":
			uv.AllowString = true
		}
	}
	return nil
}
