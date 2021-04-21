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
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
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
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
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
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
	}
}
