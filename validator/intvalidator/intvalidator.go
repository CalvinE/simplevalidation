package intvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

var (
	numberMinValueErrorTemplate = "min: the field %s value %i is less than the minimum value %i"
	numberMaxValueErrorTemplate = "max: The field %s value %i is greater than the maximum value %i"
)

// int validator

type intValidator struct {
	Min *int64
	Max *int64
}

func New() validator.Validator {
	numValidator := intValidator{}
	return &numValidator
}

func validateInt(i int64, min, max *int64, name string) (bool, error) {
	if min != nil {
		if i < *min {
			errorMessage := fmt.Sprintf(numberMinValueErrorTemplate, name, i, *min)
			return false, errors.New(errorMessage)
		}
	}
	if max != nil {
		if i > *max {
			errorMessage := fmt.Sprintf(numberMaxValueErrorTemplate, name, i, *max)
			return false, errors.New(errorMessage)
		}
	}
	return true, nil
}

func (nv *intValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	var value int64
	switch t := n.(type) {
	case int:
		value = int64(t)
	case int8:
		value = int64(t)
	case int16:
		value = int64(t)
	case int32:
		value = int64(t)
	default: // in default case we do an assertion... it will panic if the data is not the right type.
		value = n.(int64)
	}
	return validateInt(value, nv.Min, nv.Max, fieldName)
}

func (nv *intValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "min":
			min, err := strconv.ParseInt(parts[1], 0, 64)
			if err != nil {
				errorString := fmt.Sprintf("intValidator tag min value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Min = &min
		case "max":
			max, err := strconv.ParseInt(parts[1], 0, 64)
			if err != nil {
				errorString := fmt.Sprintf("intValidator tag max value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Max = &max
		}
	}
	return nil
}
