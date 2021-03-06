package uintvalidator

import (
	"reflect"
	"strings"
	"testing"
)

func TestValidUintValue(t *testing.T) {
	min, max := uint64(1), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUintValueNoMin(t *testing.T) {
	max := uint64(3)
	tValidator := uintValidator{
		Max: &max,
	}
	testValue := uint(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUintValueNoMax(t *testing.T) {
	min := uint64(1)
	tValidator := uintValidator{
		Min: &min,
	}
	testValue := uint(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUint8Value(t *testing.T) {
	min, max := uint64(1), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint8(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUint16Value(t *testing.T) {
	min, max := uint64(1), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint16(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUint32Value(t *testing.T) {
	min, max := uint64(1), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint32(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidInt64Value(t *testing.T) {
	min, max := uint64(1), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint64(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidMinInclusive(t *testing.T) {
	min, max := uint64(1), uint64(2)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint64(1)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidMaxInclusive(t *testing.T) {
	min, max := uint64(1), uint64(2)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint64(2)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidInvalidMin(t *testing.T) {
	min, max := uint64(2), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint64(1)
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
	min, max := uint64(1), uint64(2)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := uint64(3)
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
	min, max := uint64(1), uint64(3)
	tValidator := uintValidator{
		Max: &max,
		Min: &min,
	}
	testValue := float32(2)
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
	// example `validate:"uint,min=7,max=14"`
	testTag := "uint,min=7,max=14"
	tagArgs := strings.Split(testTag, ",")
	tValidator := uintValidator{}
	tValidator.ReadOptionsFromTagItems(tagArgs[1:])
	if tValidator.Min == nil || *tValidator.Min != 7 {
		t.Error("*Min should be -7")
	}
	if tValidator.Max == nil || *tValidator.Max != 14 {
		t.Error("*Max should be 7")
	}
}
