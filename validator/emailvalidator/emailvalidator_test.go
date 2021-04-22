package emailvalidator

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	PERFORM_MX_TESTS_ENV_KEY = "SIMPLEVALIDATION_PERFORM_MX_TESTS"
)

func TestValidEmail(t *testing.T) {
	tValidator := emailValidator{
		Required: true,
	}
	testValue := "user@domain.com"
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because a valid email was provided: ", err.Error())
	}
}

func TestValidNotRequiredEmptyEmail(t *testing.T) {
	tValidator := emailValidator{
		Required: false,
	}
	testValue := ""
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	if !isValid {
		t.Error("isValid should be true because Required is false and an empty string was provided: ", err.Error())
	}
}

func TestValidDomainMX(t *testing.T) {
	// os.Setenv(PERFORM_MX_TESTS_ENV_KEY, "1")
	doTest := os.Getenv(PERFORM_MX_TESTS_ENV_KEY)
	if doTest == "1" {
		tValidator := emailValidator{
			CheckDomainMX: true,
		}
		testValue := "user@gmail.com"
		valueKind := reflect.TypeOf(testValue).Kind()
		isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
		if !isValid {
			t.Error("isValid should be true because CheckDomainMX is true and the value provided has a valid domain mx record.: ", err.Error())
		}
	} else {
		t.Skipf("Skipping MX test becuase env variable %s is %s, but must be 1 to run the doamin mx validation tests", PERFORM_MX_TESTS_ENV_KEY, doTest)
	}
}

func TestInvalidRequired(t *testing.T) {
	tValidator := emailValidator{
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

func TestInvalidEmail(t *testing.T) {
	tValidator := emailValidator{}
	testValue := "not an email..."
	valueKind := reflect.TypeOf(testValue).Kind()
	isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
	expectedErrorPrefix := "invalid:"
	if isValid {
		t.Error("isValid should be false because the value provided is not a valid email adresss.")
	} else if err == nil {
		t.Error("err should be populated")
	} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
		t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
	}
}

// This test is weird, we need to find a somain that will never have an mx record, or procure one...

// func TestInvalidNoDomainMXRecord(t *testing.T) {
// 	os.Setenv(PERFORM_MX_TESTS_ENV_KEY, "1")
// 	doTest := os.Getenv(PERFORM_MX_TESTS_ENV_KEY)
// 	if doTest == "1" {
// 		tValidator := emailValidator{
// 			CheckDomainMX: true,
// 		}
// 		testValue := "user@adaqs.com"
// 		valueKind := reflect.TypeOf(testValue).Kind()
// 		isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
// 		expectedErrorPrefix := "mx missing:"
// 		if isValid {
// 			t.Error("isValid should be false because the value provided is not a valid email adresss.")
// 		} else if err == nil {
// 			t.Error("err should be populated")
// 		} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
// 			t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
// 		}
// 	} else {
// 		t.Skipf("Skipping MX test becuase env variable %s is %s, but must be 1 to run the doamin mx validation tests", PERFORM_MX_TESTS_ENV_KEY, doTest)
// 	}
// }

func TestInvalidDomainMXErrorRecord(t *testing.T) {
	// os.Setenv(PERFORM_MX_TESTS_ENV_KEY, "1")
	doTest := os.Getenv(PERFORM_MX_TESTS_ENV_KEY)
	if doTest == "1" {
		tValidator := emailValidator{
			CheckDomainMX: true,
		}
		testValue := "user@notarealdomain943292401.edu"
		valueKind := reflect.TypeOf(testValue).Kind()
		isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
		expectedErrorPrefix := "mx error:"
		if isValid {
			t.Error("isValid should be false because CheckDomainMX is true and the value provided is an email address with a non existant doamin")
		} else if err == nil {
			t.Error("err should be populated")
		} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
			t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
		}
	} else {
		t.Skipf("Skipping MX test becuase env variable %s is %s, but must be 1 to run the doamin mx validation tests", PERFORM_MX_TESTS_ENV_KEY, doTest)
	}
}

// TODO: How to simulate a failure to lookup an MX record, not that there is no mx record?
// func TestInvalidDomainLookupError(t *testing.T) {
// 	doTest := os.Getenv(PERFORM_MX_TESTS_ENV_KEY)
// 	if doTest == "1" {
// 		tValidator := emailValidator{}
// 		testValue := "not an email..."
// 		valueKind := reflect.TypeOf(testValue).Kind()
// 		isValid, err := tValidator.Validate(testValue, "testValue", valueKind)
// 		expectedErrorPrefix := "invalid:"
// 		if isValid {
// 			t.Error("isValid should be false because the value provided is not a valid email adresss.")
// 		} else if err == nil {
// 			t.Error("err should be populated")
// 		} else if hasRightIndex := strings.Index(err.Error(), expectedErrorPrefix); hasRightIndex == -1 {
// 			t.Errorf("err.Error() begin with '%s': %s", expectedErrorPrefix, err.Error())
// 		}
// 	} else {
// 		t.Skipf("Skipping MX test becuase env variable %s is %s, but must be 1 to run the doamin mx validation tests", PERFORM_MX_TESTS_ENV_KEY, doTest)
// 	}
// }

func TestInvalidType(t *testing.T) {
	tValidator := emailValidator{
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
	// example `validate:"email,required,checkdomainmx"`
	testTag := "email,required,checkdomainmx"
	tagArgs := strings.Split(testTag, ",")
	tValidator := emailValidator{}
	tValidator.ReadOptionsFromTagItems(tagArgs[1:])
	if tValidator.Required != true {
		t.Error("Required should be true")
	}
	if tValidator.CheckDomainMX != true {
		t.Error("CheckDomainMX should be true")
	}
}
