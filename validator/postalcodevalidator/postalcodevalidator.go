package postalcodevalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

// TODO: check a collection ov valid zip codes or an external service?
type postalcodeValidator struct {
	Required bool
}

const (
	postalcodeInvalidErrorTemplate  = "invalid: the field %s does cont contain a valid email. '%s' was provided"
	postalcodeRequiredErrorTemplate = "required: the field %s is required"
)

var (
	postalCodeValidationRegexp = regexp.MustCompile("^[0-9]{5}$")
)

func New() validator.Validator {
	validator := &postalcodeValidator{}
	return validator
}

func (pcv *postalcodeValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	value := n.(string)
	valueProvided := value != ""

	if pcv.Required && !valueProvided {
		errorMessage := fmt.Sprintf(postalcodeRequiredErrorTemplate, fieldName)
		return false, errors.New(errorMessage)
	}
	if valueProvided {
		if !postalCodeValidationRegexp.Match([]byte(value)) {
			errorMessage := fmt.Sprintf(postalcodeInvalidErrorTemplate, fieldName, value)
			return false, errors.New(errorMessage)
		}
	}
	return true, nil
}

func (pcv *postalcodeValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "required":
			pcv.Required = true
		}
	}
	return nil
}
