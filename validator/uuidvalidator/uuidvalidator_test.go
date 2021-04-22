package uuidvalidator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestValidUUID(t *testing.T) {
	tValidator := uuidValidator{
		Required: true,
	}
	testValue := uuid.New()
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because a valid uuid is provided: ", err.Error())
	}
}

func TestValidEmptyUUID(t *testing.T) {
	tValidator := uuidValidator{
		AllowEmptyUUID: true,
		Required:       true,
	}
	testValue := uuid.UUID{}
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because AllowEmptyUUID is true and an empty uuid was provided: ", err.Error())
	}
}

func TestValidUUIDString(t *testing.T) {
	tValidator := uuidValidator{
		AllowString: true,
		Required:    false,
	}
	testValue := uuid.New().String()
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true AllowString is true and a valid uuid string was provided: ", err.Error())
	}
}

func TestValidEmptyUUIDString(t *testing.T) {
	tValidator := uuidValidator{
		AllowEmptyUUID: true,
		AllowString:    true,
		Required:       true,
	}
	testValue := uuid.UUID{}.String()
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because AllowEmptyUUID and AllowString are true, and an empty uuid string is provided: ", err.Error())
	}
}

func TestValidNotRerquiedEmptyString(t *testing.T) {
	tValidator := uuidValidator{
		AllowString: true,
		Required:    false,
	}
	testValue := ""
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Required is false and AllowString is true, and an empty string is provided: ", err.Error())
	}
}

func TestInvalidAllowString(t *testing.T) {
	tValidator := uuidValidator{
		AllowString: false,
		Required:    true,
	}
	testValue := uuid.New().String()
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "type:"
	if isValid {
		t.Error("isValid should be false because AllowString is false and the value provided is a string")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidAllowEmptyUUID(t *testing.T) {
	tValidator := uuidValidator{
		AllowEmptyUUID: false,
		Required:       true,
	}
	testValue := uuid.UUID{}
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "no empty:"
	if isValid {
		t.Error("isValid should be false because AllowEmptyUUID is false and the value provided is an empty uuid")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidRequiredAllowString(t *testing.T) {
	tValidator := uuidValidator{
		AllowString: true,
		Required:    true,
	}
	testValue := ""
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "required:"
	if isValid {
		t.Error("isValid should be false because the value being validated is an empty string")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidStringUUIDFormat(t *testing.T) {
	tValidator := uuidValidator{
		AllowString: true,
		Required:    true,
	}
	testValue := "not a uuid..."
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "invalid:"
	if isValid {
		t.Error("isValid should be false because the value being validated is not a valid")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

func TestInvalidType(t *testing.T) {
	tValidator := uuidValidator{
		AllowString: true,
		Required:    true,
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
	// example `validate:"uuid,required,allowemptyuuid,allowstring"`
	testTag := "uuid,required,allowemptyuuid,allowstring"
	tagArgs := strings.Split(testTag, ",")
	tValidator := uuidValidator{}
	tValidator.ReadOptionsFromTagItems(tagArgs[1:])
	if tValidator.Required != true {
		t.Error("Required should be true")
	}
	if tValidator.AllowEmptyUUID != true {
		t.Error("AllowEmptyUUID should be true")
	}
	if tValidator.AllowString != true {
		t.Error("AllowString should be true")
	}
}
