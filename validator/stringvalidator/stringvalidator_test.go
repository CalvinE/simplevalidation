package stringvalidator

import (
	"reflect"
	"strings"
	"testing"
)

func TestValidString(t *testing.T) {
	min, max := 8, 16
	tValidator := stringValidator{
		Min:      &min,
		Max:      &max,
		Required: true,
	}
	testValue := "this is a string"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because a valid string is provided: ", err.Error())
	}
}

func TestValidMinLength(t *testing.T) {
	min := 3
	tValidator := stringValidator{
		Min:      &min,
		Required: true,
	}
	testValue := "hello"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Min is less than the provided values length: ", err.Error())
	}
}

func TestValidMaxLength(t *testing.T) {
	max := 8
	tValidator := stringValidator{
		Max:      &max,
		Required: true,
	}
	testValue := "hello"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Max is greater than the provided values length: ", err.Error())
	}
}

func TestValidMinLengthInclusive(t *testing.T) {
	min := 3
	tValidator := stringValidator{
		Min:      &min,
		Required: true,
	}
	testValue := "hey"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Min is equal to the provided values length: ", err.Error())
	}
}

func TestValidMaxLengthInclusive(t *testing.T) {
	max := 5
	tValidator := stringValidator{
		Max:      &max,
		Required: true,
	}
	testValue := "hello"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Max is equal to the provided values length: ", err.Error())
	}
}

func TestValidNotRequiredEmptyString(t *testing.T) {
	tValidator := stringValidator{
		Required: false,
	}
	testValue := ""
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Required is false, and an empty string is provided: ", err.Error())
	}
}

func TestInvalidRequired(t *testing.T) {
	tValidator := stringValidator{
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

func TestInvalidMin(t *testing.T) {
	min := 5
	tValidator := stringValidator{
		Min:      &min,
		Required: true,
	}
	testValue := "test"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "min length:"
	if isValid {
		t.Error("isValid should be false because Min is greater than the length of the value provided")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidMax(t *testing.T) {
	max := 5
	tValidator := stringValidator{
		Max:      &max,
		Required: true,
	}
	testValue := "testing"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "max length:"
	if isValid {
		t.Error("isValid should be false because Max is less than the length of the value provided")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidType(t *testing.T) {
	tValidator := stringValidator{
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
	// example `validate:"string,required,min=3,max=8"`
	testTag := "string,required,min=3,max=8"
	tagArgs := strings.Split(testTag, ",")
	tValidator := stringValidator{}
	tValidator.ReadOptionsFromTagItems(tagArgs[1:])
	if tValidator.Required != true {
		t.Error("Required should be true")
	}
	if tValidator.Min == nil || *tValidator.Min != 3 {
		t.Error("*Min should be 3")
	}
	if tValidator.Max == nil || *tValidator.Max != 8 {
		t.Error("*Max should be 8")
	}
}
