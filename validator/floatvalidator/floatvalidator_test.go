package floatvalidator

import (
	"reflect"
	"strings"
	"testing"
)

func TestValidFloat32Value(t *testing.T) {
	min, max := float64(1.1), float64(3.3)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float32(2.2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidFloat64Value(t *testing.T) {
	min, max := float64(1.1), float64(3.3)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float64(2.2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidFloatValueNoMin(t *testing.T) {
	max := float64(3)
	tValidator := floatValidator{
		Max: &max,
	}
	testValue := float64(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidFloatValueNoMax(t *testing.T) {
	min := float64(1)
	tValidator := floatValidator{
		Min: &min,
	}
	testValue := float64(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidMinInclusive(t *testing.T) {
	min, max := float64(1.1), float64(2.2)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float64(1.1)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidMaxInclusive(t *testing.T) {
	min, max := float64(1.1), float64(2.2)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float64(2.2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidInvalidMin(t *testing.T) {
	min, max := float64(1.2), float64(3.3)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float64(1.1)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "min:"
	if isValid {
		t.Error("isValid should be false because the the value being validated is less than the Min value")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestValidInvalidMax(t *testing.T) {
	min, max := float64(1.1), float64(3.3)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float64(3.4)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "max:"
	if isValid {
		t.Error("isValid should be false because the the value being validated is more than the Max value")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestValidInvalidType(t *testing.T) {
	min, max := float64(1.1), float64(3.3)
	tValidator := floatValidator{
		Max: &max,
		Min: &min,
	}
	testValue := int(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "type:"
	if isValid {
		t.Error("isValid should be false because the value being validated is not a valid type")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestReadOptionsFromTagItemsAllParameters(t *testing.T) {
	// example `validate:"float,min=-7.7,max=7.7"`
	testTag := "float,min=-7.7,max=7.7"
	tagArgs := strings.Split(testTag, ",")
	tValidator := floatValidator{}
	tValidator.ReadOptionsFromTagItems(tagArgs[1:])
	if tValidator.Min == nil || *tValidator.Min != -7.7 {
		t.Error("*Min should be -7")
	}
	if tValidator.Max == nil || *tValidator.Max != 7.7 {
		t.Error("*Max should be 7")
	}
}
