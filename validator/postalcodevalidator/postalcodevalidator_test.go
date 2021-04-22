package postalcodevalidator

import (
	"reflect"
	"strings"
	"testing"
)

func TestValidPostalCode(t *testing.T) {
	tValidator := postalcodeValidator{
		Required: true,
	}
	testValue := "31008"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because a valid postal code was provided: ", err.Error())
	}
}

func TestValidNotRequiredEmptyPostalCode(t *testing.T) {
	tValidator := postalcodeValidator{
		Required: false,
	}
	testValue := ""
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Required is false and an empty string was provided: ", err.Error())
	}
}

func TestInvalidRequired(t *testing.T) {
	tValidator := postalcodeValidator{
		Required: true,
	}
	testValue := ""
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "required:"
	if isValid {
		t.Error("isValid should be false because Required is true and the value provided is an empty string")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidPostalCode(t *testing.T) {
	tValidator := postalcodeValidator{}
	testValue := "not an postal code..."
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "invalid:"
	if isValid {
		t.Error("isValid should be false because the value provided is not a valid postal code.")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidType(t *testing.T) {
	tValidator := postalcodeValidator{
		Required: true,
	}
	testValue := uint(32)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "type:"
	if isValid {
		t.Error("isValid should be false because the value being validated is not of a valid type.")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestReadOptionsFromTagItemsAllParameters(t *testing.T) {
	// example `validate:"postalcode,required"`
	testTag := "postalcode,required"
	tagArgs := strings.Split(testTag, ",")
	tValidator := postalcodeValidator{}
	tValidator.ReadOptionsFromTagItems(tagArgs[1:])
	if tValidator.Required != true {
		t.Error("Required should be true")
	}
}
