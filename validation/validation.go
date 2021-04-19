package validation

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/calvine/simplevalidation/validation/validationparams"
	"github.com/calvine/simplevalidation/validator"
	"github.com/calvine/simplevalidation/validator/emailvalidator"
	"github.com/calvine/simplevalidation/validator/floatvalidator"
	"github.com/calvine/simplevalidation/validator/intvalidator"
	"github.com/calvine/simplevalidation/validator/postalcodevalidator"
	"github.com/calvine/simplevalidation/validator/stringvalidator"
	"github.com/calvine/simplevalidation/validator/uintvalidator"
	"github.com/calvine/simplevalidation/validator/uuidvalidator"
)

type ValidationError struct {
	DataType string
	Errors   validationErrorMap
}

type validationErrorMap map[string][]error

func (e *ValidationError) Error() string {
	var errorBuffer bytes.Buffer
	headerLine := fmt.Sprintf("%s object validation failed:", e.DataType)
	fmt.Fprint(&errorBuffer, headerLine)
	for key, error := range e.Errors {
		fmt.Fprintf(&errorBuffer, "\t%s: %s", key, error)
	}
	return errorBuffer.String()
}

var (
	pointerNilTemplate string                                = "field %s was nil but is required"
	validators         map[string]validator.ValidatorFactory = map[string]validator.ValidatorFactory{
		"email":      emailvalidator.New,
		"float":      floatvalidator.New,
		"int":        intvalidator.New,
		"uint":       uintvalidator.New,
		"postalcode": postalcodevalidator.New,
		"string":     stringvalidator.New,
		"uuid":       uuidvalidator.New,
	}
)

func getValidatorFromTag(validatorName, fieldName string) (validator.Validator, error) {
	if validatorName == "struct" {
		return nil, nil
	}
	typeValidatorFactory, ok := validators[validatorName]
	if !ok {
		errMsg := fmt.Sprintf("Validator of type %s is not registered.", validatorName)
		return nil, errors.New(errMsg)
	}
	return typeValidatorFactory(), nil
}

func getValidatorInfo(validatorName string) (name string, arrayDepth uint8) {
	arrayDepth = 0
	tagNameStartIndex := 0
	if len(validatorName) > 2 {
		for {
			if validatorName[tagNameStartIndex] == '[' && validatorName[tagNameStartIndex+1] == ']' {
				arrayDepth++
				tagNameStartIndex += 2
			} else {
				break
			}
		}
	}
	return validatorName[tagNameStartIndex:], arrayDepth
}

func performFieldValidation(validationInfo validationparams.ValidationParams, validationErrors *validationErrorMap) {
	fieldErrors := []error{}
	value := reflect.ValueOf(validationInfo.Value)
	kind := value.Kind()
	vType := value.Type()
	// fmt.Printf("%v - %v\n\n", kind, vType.String())
	switch kind {
	case reflect.Ptr:
		isNil := value.IsNil()
		if !validationInfo.Required && isNil {
			// If there is only 1 tag arg, its the validation type.
			// If required is not the second arg then its not required.
			return
		} else if validationInfo.Required && isNil {
			errorMessage := fmt.Sprintf(pointerNilTemplate, validationInfo.Name)
			fieldErrors = append(fieldErrors, errors.New(errorMessage))
		} else {
			fieldValue := value.Elem().Interface()
			recursiveFieldValidator := validationparams.ValidationParams{
				ArrayDepth:     validationInfo.ArrayDepth,
				Name:           validationInfo.Name,
				FieldValidator: validationInfo.FieldValidator,
				Required:       validationInfo.Required,
				StructDepth:    validationInfo.StructDepth,
				Value:          fieldValue,
			}
			performFieldValidation(recursiveFieldValidator, validationErrors)
		}
	case reflect.Struct:
		// handle embedded structs...
		structDepth := validationInfo.StructDepth + 1
		for i := 0; i < value.NumField(); i++ {
			field := vType.Field(i)
			tag := field.Tag.Get("validate")
			// fieldKind := field.Type.Kind()
			// fieldType := field.Type
			// fmt.Printf("k: %v - t: %v\n\n", fieldKind, fieldType)
			if tag == "" || tag == "-" {
				continue
			}
			fieldName := field.Name
			if structDepth > 1 {
				// this is specifically for structs within structs to create a better reference to the name of the field being validated.
				fieldName = fmt.Sprintf("%s.%s", validationInfo.Name, fieldName)
			}
			fieldValue := value.Field(i)
			tagArgs := strings.Split(tag, ",")
			required := len(tagArgs) > 1 && tagArgs[1] == "required"
			validatorName, arrayDepth := getValidatorInfo(tagArgs[0])
			validator, err := getValidatorFromTag(validatorName, fieldName)
			if err != nil {
				// make a custom type not registered error?
				fieldErrors = append(fieldErrors, err)
				// should we panic?
			} else if validator == nil {
				validationData := validationparams.ValidationParams{
					ArrayDepth:     arrayDepth,
					FieldValidator: validator,
					Name:           fieldName,
					Required:       required,
					StructDepth:    structDepth,
					Value:          fieldValue.Interface(),
				}
				performFieldValidation(validationData, validationErrors)
			} else {
				err = validator.ReadOptionsFromTagItems(tagArgs[1:])
				if err != nil {
					// make a custom tag invalid error?
					fieldErrors = append(fieldErrors, err)
				}
				validationData := validationparams.ValidationParams{
					ArrayDepth:     arrayDepth,
					FieldValidator: validator,
					Name:           fieldName,
					Required:       required,
					StructDepth:    structDepth,
					Value:          fieldValue.Interface(),
				}
				performFieldValidation(validationData, validationErrors)
			}
		}
	default:
		// perform normal field validation
		if validationInfo.ArrayDepth == 0 {
			_, fieldError := validationInfo.FieldValidator.Validate(validationInfo.Value, validationInfo.Name, kind)
			if fieldError != nil {
				fieldErrors = append(fieldErrors, fieldError)
			}
		} else {
			// potentially nested array element validation.
			currentArrayDepth := validationInfo.ArrayDepth
			switch reflect.TypeOf(validationInfo.Value).Kind() {
			case reflect.Slice, reflect.Array:
				currentArrayDepth--
				currentLevelSlice := reflect.ValueOf(validationInfo.Value)
				for i := 0; i < currentLevelSlice.Len(); i++ {
					performFieldValidation(validationparams.ValidationParams{
						ArrayDepth:     currentArrayDepth,
						FieldValidator: validationInfo.FieldValidator,
						Name:           fmt.Sprintf("%s[%d]", validationInfo.Name, i),
						Required:       validationInfo.Required,
						StructDepth:    validationInfo.StructDepth,
						Value:          currentLevelSlice.Index(i).Interface(),
					}, validationErrors)
				}
			default:
				// This should not happen. add error...
			}
		}
	}
	if len(fieldErrors) > 0 {
		(*validationErrors)[validationInfo.Name] = fieldErrors
	}
}

// Thinking of building validation code to read tags.
// The Validator parameter is present to allow for validating non struct values. In this event A *Validator can be passed in and evaluated on a non struct value link an individual int or string
// TODO: Validate nested structs and arrays.
func Validate(v *validationparams.ValidationParams) (*ValidationError, error) {
	validationErrors := validationErrorMap{}
	if v == nil {
		return nil, errors.New("no FieldValidationData provided")
	}
	performFieldValidation(*v, &validationErrors)
	if len(validationErrors) > 0 {
		return &ValidationError{
			Errors: validationErrors,
		}, nil
	}
	return nil, nil
}

func ValidateStructWithTag(s interface{}) *ValidationError {
	validationErrors := validationErrorMap{}
	validationData := validationparams.New()
	validationData.Value = s
	performFieldValidation(validationData, &validationErrors)
	if len(validationErrors) > 0 {
		return &ValidationError{
			Errors: validationErrors,
		}
	}
	return nil
}

func RegisterValidator(name string, validator func() validator.Validator) {
	validators[name] = validator
}