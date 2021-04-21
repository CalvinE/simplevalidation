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
	AllowEmptyUUID bool
	AllowString    bool
	Required       bool
}

const (
	uuidStringProvidedErrorTemplate = "type: The field %s is a string, but allowstring was not provided in the validation tag."
	uuidInvalidStringErrorTemplate  = "invalid: the field %s has the value %s which could not be parsed into a UUID."
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
	noValueProvided := false
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
		}
	default:
		var ok bool
		value, ok = t.(uuid.UUID)
		if !ok {
			errorMessage := fmt.Sprintf(validator.InvalidTypeErrorTemplate, fieldName, t)
			return false, errors.New(errorMessage)
		}
		if (t == uuid.UUID{}) {
			noValueProvided = !uv.AllowEmptyUUID
		}
	}

	if uv.Required && noValueProvided {
		errorMessage := fmt.Sprintf(uuidNoValueErrorTemplate, fieldName, value)
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
