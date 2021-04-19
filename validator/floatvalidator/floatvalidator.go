package floatvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

var (
	numberMinValueErrorTemplate = "The field %s value %i is less than the minimum value %i"
	numberMaxValueErrorTemplate = "The field %s value %i is greater than the maximum value %i"
)

type floatValidator struct {
	Min *float64
	Max *float64
}

func New() validator.Validator {
	numValidator := floatValidator{}
	return &numValidator
}

func validateFloat(i float64, min, max *float64, name string) (bool, error) {
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

func (nv *floatValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	var value float64
	switch t := n.(type) {
	case float32:
		value = float64(t)
	default: // in default case we do an assertion... it will panic if the data is not the right type.
		value = n.(float64)
	}
	return validateFloat(value, nv.Min, nv.Max, fieldName)
}

func (nv *floatValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "min":
			min, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				errorString := fmt.Sprintf("floatValidator tag min value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Min = &min
		case "max":
			max, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				errorString := fmt.Sprintf("floatValidator tag max value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Max = &max
		}
	}
	return nil
}
