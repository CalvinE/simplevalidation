package validation

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// This is for testing embedded structs
type OtherThing struct {
	ID          int    `validate:"int,required,min=3"`
	Description string `validate:"string,min=15,max=150"`
}

// this is for testing fields with struct values
type Details struct {
	Name  string `validate:"string,max=50"`
	Value int    `validate:"int,min=0,max=150"`
}

type SimpleItem struct {
	Name        string  `validate:"string,required"`
	Description *string `validate:"string,required"`
	Other       string  `validate:"-"`
	Other2      string  `validate:""`
	Other3      string  ``
}

type TestStruct struct {
	Arry       *[]int     `validate:"[]int,min=0,max=150"`
	ID         uuid.UUID  `validate:"uuid,required,allowemptyuuid"`
	ID2        *uuid.UUID `validate:"uuid"`
	Age        int        `validate:"int,min=0,max=150"`
	Email      string     `validate:"email,required"`
	Name       string     `validate:"string,required,min=3,max=50"`
	PostalCode string     `validate:"postalcode,required"`
	Score      *uint16    `validate:"uint,required,min=1,max=1500"`
	TheTime    time.Time  `validate:"time,required"`
	Detail     Details    `validate:"struct"`
	OtherThing `validate:"struct"`
}

// TOOD: Build tests for embedded structs and arrays of structs.
func TestTagValidation(t *testing.T) {
	var score uint16 = 7
	var arryData = []int{1, 2, 3, 4, 5}
	testStruct := TestStruct{
		Age:        33,
		Arry:       &arryData,
		Email:      "test@user.com",
		Name:       "Calvin",
		PostalCode: "32105",
		Score:      &score,
		OtherThing: OtherThing{
			ID: 11,
		},
		TheTime: time.Now(),
	}
	validationError := ValidateStructWithTag(testStruct)
	// TODO: almost done errors array is weird...
	if validationError != nil {
		t.Error("Oops.")
	}
}

func TestRequireNilPointer(t *testing.T) {
	testValue := SimpleItem{
		Name:        "thing",
		Description: nil,
	}
	validationError := ValidateStructWithTag(testValue)
	if validationError == nil {
		t.Error("Description should have cause an error because its required, but a nil pointer.")
	} else if _, ok := validationError.Errors["Description"]; !ok {
		t.Error("validationError.Errors should contain key 'Description")
	}
}

func TestNoOrSkipTag(t *testing.T) {
	s := "test"
	testValue := SimpleItem{
		Name:        "thing",
		Description: &s,
	}
	validationError := ValidateStructWithTag(testValue)
	if validationError != nil {
		t.Error("testValue is valid and should not have resulted in an error", validationError.Error())
	}
}

func TestNoValidatorRegistered(t *testing.T) {
	testValue := struct {
		Name string `validate:"notarealvalidator"`
	}{
		"not important",
	}
	validationError := ValidateStructWithTag(testValue)
	if validationError == nil {
		t.Error("Name should have an error in the error map")
	} else {
		errorValue, ok := validationError.Errors["value"]
		if !ok {
			t.Error("Name should be in the validationError.Errors map: ", validationError.Error())
		} else if hasExpectedError := strings.Index(errorValue[0].Error(), "no validator:"); hasExpectedError == -1 {
			t.Error("the only error in under Name should be the no validator error.", errorValue[0].Error())
		}

	}
}
