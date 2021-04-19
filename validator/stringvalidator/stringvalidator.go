package stringvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

// Struct that contains the fields required to validate a string.
type stringValidator struct {
	// The min length allowed for the input string.
	Min *int
	// The max length allowed for the input string.
	Max *int
	// If true then an empty string ("") will cause a validation error.
	Required bool
	//pattern string
}

var (
	stringMinLengthTemplate = "The value of %s is %s of length %d which is less than the minimum length %d"
	stringMaxLengthTemplate = "The value of %s is %s of length %d which is greater than the maximum length %d"
	stringRequiredTemplate  = "The value of %s is blank"
)

func New() validator.Validator {
	strValidator := stringValidator{}
	return &strValidator
}

func (nv *stringValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	stringValue := n.(string)
	valueLength := len(stringValue)
	if nv.Min != nil {
		if valueLength < *nv.Min {
			errorMessage := fmt.Sprintf(stringMinLengthTemplate, fieldName, valueLength, stringValue, *nv.Min)
			return false, errors.New(errorMessage)
		}
	}
	if nv.Max != nil {
		if valueLength > *nv.Max {
			errorMessage := fmt.Sprintf(stringMaxLengthTemplate, fieldName, valueLength, stringValue, *nv.Min)
			return false, errors.New(errorMessage)
		}
	}
	if nv.Required && valueLength == 0 {
		errorMessage := fmt.Sprintf(stringRequiredTemplate, fieldName)
		return false, errors.New(errorMessage)
	}
	return true, nil

}

func (nv *stringValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "min":
			min, err := strconv.Atoi(parts[1])
			if err != nil {
				errorString := fmt.Sprintf("stringValidator tag min value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Min = &min
		case "max":
			max, err := strconv.Atoi(parts[1])
			if err != nil {
				errorString := fmt.Sprintf("stringValidator tag max value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Max = &max
		case "required":
			nv.Required = true
		}
	}
	return nil
}
