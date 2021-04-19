package emailvalidator

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/calvine/simplevalidation/validator"
)

// TODO: add a domain whitelist so the user can specify valid domains in the tag? seperated by ';'?
type emailValidator struct {
	// If true we check for a valid MX record for the domain of the email.
	CheckDomainMX bool
	Required      bool
}

var (
	emailInvalidErrorTemplate          = "The field %s does cont contain a valid email. '%s' was provided"
	emailRequiredErrorTemplate         = "The field %s is required"
	emailDomainMXNotFoundErrorTemplate = "The field %s had no MX records found for domain %s"
	emailDomainMXLookupErrorTemplate   = "The field %s encountered an error occurred while validating domain MX record for domain %s. Error: %s"

	emailValidationRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func New() validator.Validator {
	validator := &emailValidator{}
	return validator
}

func (ev *emailValidator) Validate(n interface{}, fieldName string, fieldKind reflect.Kind) (bool, error) {
	value := n.(string)
	valueProvided := value != ""

	if ev.Required && !valueProvided {
		errorMessage := fmt.Sprintf(emailRequiredErrorTemplate, fieldName)
		return false, errors.New(errorMessage)
	}
	if valueProvided {
		if !emailValidationRegexp.Match([]byte(value)) {
			errorMessage := fmt.Sprintf(emailInvalidErrorTemplate, fieldName, value)
			return false, errors.New(errorMessage)
		} else if ev.CheckDomainMX {
			emailParts := strings.Split(value, "@")
			domain := emailParts[1]
			mx, err := net.LookupMX(emailParts[1])
			if err != nil {
				errorMessage := fmt.Sprintf(emailDomainMXLookupErrorTemplate, fieldName, domain, err.Error())
				return false, errors.New(errorMessage)
			} else if len(mx) == 0 {
				errorMessage := fmt.Sprintf(emailDomainMXNotFoundErrorTemplate, fieldName, domain)
				return false, errors.New(errorMessage)
			}
		}
	}
	return true, nil
}

func (ev *emailValidator) ReadOptionsFromTagItems(items []string) error {
	for i := 0; i < len(items); i++ {
		switch parts := strings.Split(items[i], "="); parts[0] {
		case "required":
			ev.Required = true
		case "checkdomainmx":
			ev.CheckDomainMX = true
		}
	}
	return nil
}
