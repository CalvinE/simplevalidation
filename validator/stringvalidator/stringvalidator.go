package stringvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

type stringValidator struct {
	Min      *int
	Max      *int
	Required bool
	//pattern string
}

var (
	stringMinLengthTemplate = "The value of %s is %s which is less than the minimum length %i"
	stringMaxLengthTemplate = "The value of %s is %s which is greater than the maximum length %i"
	stringRequiredTemplate  = "The value of %s is blank"
)

func New() validator.Validator {
	strValidator := stringValidator{}
	return &strValidator
}

func (nv *stringValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	stringValue := n.(string)
	if nv.Min != nil {
		if len(stringValue) < *nv.Min {
			errorMessage := fmt.Sprintf(stringMinLengthTemplate, fieldName, stringValue, *nv.Min)
			return false, errors.New(errorMessage)
		}
	}
	if nv.Max != nil {
		if len(stringValue) > *nv.Max {
			errorMessage := fmt.Sprintf(stringMaxLengthTemplate, fieldName, stringValue, *nv.Min)
			return false, errors.New(errorMessage)
		}
	}
	if nv.Required && len(stringValue) == 0 {
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
