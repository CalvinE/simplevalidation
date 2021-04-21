package uintvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

const (
	numberMinValueErrorTemplate = "min: the field %s value %d is less than the minimum value %d"
	numberMaxValueErrorTemplate = "max: the field %s value %d is greater than the maximum value %d"
)

// uint validator

type uintValidator struct {
	Min *uint64
	Max *uint64
}

func New() validator.Validator {
	numValidator := uintValidator{}
	return &numValidator
}

func validateUint(i uint64, min, max *uint64, name string) (bool, error) {
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

func (nv *uintValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	var value uint64
	switch t := n.(type) {
	case uint:
		value = uint64(t)
	case uint8:
		value = uint64(t)
	case uint16:
		value = uint64(t)
	case uint32:
		value = uint64(t)
	default: // in default case we do an assertion... it will panic if the data is not the right type.
		var ok bool
		value, ok = n.(uint64)
		if !ok {
			errorMessage := fmt.Sprintf(validator.InvalidTypeErrorTemplate, fieldName, t)
			return false, errors.New(errorMessage)
		}
	}
	return validateUint(value, nv.Min, nv.Max, fieldName)
}

func (nv *uintValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "min":
			min, err := strconv.ParseUint(parts[1], 0, 64)
			if err != nil {
				errorString := fmt.Sprintf("uintValidator tag min value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Min = &min
		case "max":
			max, err := strconv.ParseUint(parts[1], 0, 64)
			if err != nil {
				errorString := fmt.Sprintf("uintValidator tag max value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			nv.Max = &max
		}
	}
	return nil
}
