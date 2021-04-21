package timevalidator

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestValidTimeValue(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := time.Unix(1618968930, 0)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidInt64Value(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968930)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUintValueNoNotBefore(t *testing.T) {
	naf := int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Required: true,
	}
	testValue := int64(1618968930)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidUintValueNoNotAfter(t *testing.T) {
	nbf := int64(1618968920)
	tValidator := timeValidator{
		AllowInt: true,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968930)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidNotBeforeInclusive(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968920)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidNotAfterInclusive(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968940)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true, but an error occurred: ", err.Error())
	}
}

func TestValidInvalidNotBefore(t *testing.T) {
	nbf, naf := int64(1618968931), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968930)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "not before:"
	if isValid {
		t.Error("isValid should be false because the value being validated is before the Not Before value")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
	}
}

func TestValidInvalidNotAfter(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968929)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968940)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "not after:"
	if isValid {
		t.Error("isValid should be false because the value being validated is after the Not After value")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
	}
}

func TestValidInvalidAllowInt(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: false,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := int64(1618968930)
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "type:"
	if isValid {
		t.Error("isValid should false be because AllowInt is false, but the validation input is an int64")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
	}
}

func TestValidInvalidRequired(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := time.Time{}
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "required:"
	if isValid {
		t.Error("isValid should be false because Required is true, and a default time value is provided")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s'", expectedErrorPrefix)
	}
}

func TestValidInvalidType(t *testing.T) {
	nbf, naf := int64(1618968920), int64(1618968940)
	tValidator := timeValidator{
		AllowInt: true,
		Naf:      &naf,
		Nbf:      &nbf,
		Required: true,
	}
	testValue := float32(1618968930)
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
