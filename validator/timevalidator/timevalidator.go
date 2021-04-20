package timevalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/calvine/simplevalidation/validator"
)

// TODO add property to pass in a layout for parsing strings so they can be validated?
type timeValidator struct {
	// AllowInt allows the user to apply this
	AllowInt bool
	// Nbf stands for Not Before and must be a unix timestamp that can be cast as an int64.
	Nbf *int64
	// Nbf stands for Not After and must be a unix timestamp that can be cast as an int64.
	Naf      *int64
	Required bool
}

const (
	timeIntNotAllowedErrorTemplate = "type: the field %s is an int, but AllowInt was false"
	timeNotBeforeErrorTemplate     = "not before: the field %s has a value of '%s' which is before '%s'"
	timeNotAfterErrorTemplate      = "not after: the field %s has a value of '%s' which is after '%s'"
	timeNoValueErrorTemplate       = "required: The field %s had an empty value: %v"
)

func New() validator.Validator {
	tValidator := timeValidator{}
	return &tValidator
}

func (tv *timeValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	var value int64
	noValueProvided := false
	switch t := n.(type) {
	case time.Time:
		value = t.Unix()
		noValueProvided = t == time.Time{}
	case int64:
		if !tv.AllowInt {
			errorMessage := fmt.Sprintf(timeIntNotAllowedErrorTemplate, fieldName)
			return false, errors.New(errorMessage)
		}
		value = t
	default:
		// var ok bool
		// value, ok = t.(int64)
		// if !ok {
		errorMessage := fmt.Sprintf(validator.InvalidTypeErrorTemplate, fieldName, t)
		return false, errors.New(errorMessage)
		// }
	}

	if tv.Nbf != nil {
		if value < *tv.Nbf {
			nbfString := time.Unix(*tv.Nbf, 0).UTC().String()
			valueString := time.Unix(value, 0).UTC().String()
			errorMessage := fmt.Sprintf(timeNotBeforeErrorTemplate, fieldName, valueString, nbfString)
			return false, errors.New(errorMessage)
		}
	}

	if tv.Naf != nil {
		if value > *tv.Naf {
			nafString := time.Unix(*tv.Naf, 0).UTC().String()
			valueString := time.Unix(value, 0).UTC().String()
			errorMessage := fmt.Sprintf(timeNotAfterErrorTemplate, fieldName, valueString, nafString)
			return false, errors.New(errorMessage)
		}
	}

	if tv.Required && noValueProvided {
		errorMessage := fmt.Sprintf(timeNoValueErrorTemplate, fieldName, value)
		return false, errors.New(errorMessage)
	}
	return true, nil

}

func (tv *timeValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "allowint":
			tv.AllowInt = true
		case "required":
			tv.Required = true
		case "nbf":
			value, err := strconv.ParseInt(parts[1], 0, 64)
			if err != nil {
				errorString := fmt.Sprintf("timeValidator tag nbf value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			tv.Nbf = &value
		case "naf":
			value, err := strconv.ParseInt(parts[1], 0, 64)
			if err != nil {
				errorString := fmt.Sprintf("timeValidator tag naf value invalid: %s", err.Error())
				return errors.New(errorString)
			}
			tv.Naf = &value
		}
	}
	return nil
}
