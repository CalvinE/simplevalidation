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

/*
	ValidationError represents the overall validation state of a value.
*/
type ValidationError struct {
	DataType string
	Errors   validationErrorMap
}

/*
	validationErrorMap is an alias for a:
		map[string][]error

	The key for the map is the field being validated. when validating a struct its the name of a field, when validting a non struct value, it defaults to "value".
*/
type validationErrorMap map[string][]error

/*
	Error produces a string that relays all of the validation failures from validation.
*/
func (e *ValidationError) Error() string {
	var errorBuffer bytes.Buffer
	dataType := e.DataType
	if dataType == "" {
		dataType = "value"
	}
	headerLine := fmt.Sprintf("%s validation failed:", dataType)
	fmt.Fprint(&errorBuffer, headerLine)
	for key, error := range e.Errors {
		fmt.Fprintf(&errorBuffer, "\t%s: %s", key, error)
	}
	return errorBuffer.String()
}

var (
	// pointerNilTemplate is the error message template then a required value is a pointer and also nil.
	pointerNilTemplate string = "field %s was nil but is required"
	// validators is a map contains a contains a validator.ValidatorFactory for each registered validator name.
	validators map[string]validator.ValidatorFactory = map[string]validator.ValidatorFactory{
		"email":      emailvalidator.New,
		"float":      floatvalidator.New,
		"int":        intvalidator.New,
		"uint":       uintvalidator.New,
		"postalcode": postalcodevalidator.New,
		"string":     stringvalidator.New,
		"uuid":       uuidvalidator.New,
	}
)

// getValidatorFromTag Takes in the validator name from the tag data and returns an instance of the appropriate validator.
// When the validatorName parameter is not registererd in the validators map, the function returns an error.
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

// getValidatorInfo reads in the raw validator name from the tag data, and parses pairs of square brackets to determin the ArrayDepth of the value being validated.
// It returns the plain validator name (with any square bracket pairs removed) for looking up in the validators map, and the array depth for the validator to use.
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

/*
	performFieldValidation is the core of the validation work flow. it takes a ValidationParams struct and a pointer to a validationErrorMap.
 	It then proceeds to call its self recursivly, until all validation is completed. Upon completion the validationErrors parameter is populated with all errors arising from validation.

	This function handles the following cases:
		- When the value being validated is a pointer it is dereferenced, and the validated.
			- When that pointer is nil validation is skipped, unless the validationparams.ValidationParams.Required field is true, then it will register a validation error.
		- When the field being validated is a struct the struct fields are traversed and the function attempts to build the appropriate validator based on the validator tag data.
		- When the field is any other kind it will attempt to validate the value, if the validationparams.ValidationParams.ArrayDepth is greater than 0 the function will iterate of the array / slice and validate each value for each level of array / slice.
*/
func performFieldValidation(validationInfo validationparams.ValidationParams, validationErrors *validationErrorMap) {
	fieldErrors := []error{}
	value := reflect.ValueOf(validationInfo.Value)
	kind := value.Kind()
	vType := value.Type()
	// fmt.Printf("%v - %v\n\n", kind, vType.String())
	if kind == reflect.Ptr {
		// handle pointers.
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
	} else if kind == reflect.Struct && validationInfo.FieldValidator == nil {
		// handle structs and embedded structs.
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
			// TODO: Add struct case as special key word?
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
	} else if validationInfo.FieldValidator != nil {
		// perform normal field validation.
		if validationInfo.ArrayDepth == 0 {
			_, fieldError := validationInfo.FieldValidator.Validate(validationInfo.Value, validationInfo.Name, kind)
			if fieldError != nil {
				fieldErrors = append(fieldErrors, fieldError)
			}
		} else {
			// potentially nested array element validation.
			currentArrayDepth := validationInfo.ArrayDepth
			// kind := reflect.TypeOf(validationInfo.Value).Kind()
			switch kind {
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
	} // else { panic? }
	if len(fieldErrors) > 0 {
		(*validationErrors)[validationInfo.Name] = fieldErrors
	}
}

//The Validator parameter is present to allow for validating non struct values. In this function A Validator pointer can be passed in and evaluated on a non struct value like an individual int or string.
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

// This function validates an input struct based on the validation tags is has in its tag data.
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

// This allows you to register custom validator to be read from struct field tag validatoin data.
func RegisterValidator(name string, customValidatorFactory validator.ValidatorFactory) {
	validators[name] = customValidatorFactory
}
