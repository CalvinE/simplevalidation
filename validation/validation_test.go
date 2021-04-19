package validation

import (
	"testing"

	"github.com/google/uuid"
)

// This is for testing embedded structs
type OtherThing struct {
	ID          int    `validate:"int,required,min=3"`
	Description string `validate:"string,min=0,max=150"`
}

// this is for testing fields with struct values
type Details struct {
	Name  string `validate:"string,max=50"`
	Value int    `validate:"int,min=0,max=150"`
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
	}
	validationError := ValidateStructWithTag(testStruct)
	// TODO: almost done errors array is weird...
	if validationError != nil {
		t.Error("Oops.")
	}
}
