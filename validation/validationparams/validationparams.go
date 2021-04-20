package validationparams

import "github.com/calvine/simplevalidation/validator"

// ValidationParams is a struct that represents the parameters for validating a value.
type ValidationParams struct {
	/*
		ArrayDepth tells the validator how many layers deep an array is.

		So for instance:

			[][]int would have an ArrayDepth of 2
			[]int would have an ArrayDepth of 1
			int would have an ArrayDepth of 0

		When utilizing tag based validation for struct fields,
		the array depth is computed by reading the validator type and counting the
		square bracket pairs before the validator name.

		While being validated the validation function will recursivly go through each arra level and validate each value.
		For any failed validation, the resulting error label will be in the format of "type[index1][index2][index3]" for as may nested arrays you may have.
	*/
	ArrayDepth uint8
	/*
		FieldValidator is the validator used to validate a value. If the validation is being performed based on struct tags, this is populated by the validation tag parser.
	*/
	FieldValidator validator.Validator
	/*
		Name is the name of the field being validated if it can be determined and if not,
		for instance if ou are validating a non struct field, you can populate it, or
	*/
	Name string
	/*
		Required is a special validation cruteria that is technically valid for any validator. If the Validator does not directly use the required flag, it is still used in the event that the value being validated is a pointer. If the value being validated is a pointer and that pointer is nil then that will result in a validation error when required is true.
	*/
	Required bool
	/*
		This is a counter that keeps track of how many structs deep validation is being performed.

		For instance:

			type s2 struct {
				C int `validation:"int"`
			}

			type s1 struct {
				A int `validation:"int,min=3,max=15"`
				B s2 `validation:"struct"`
			}

		When validating A StructDepth is 1 and when validating C StructDepth is 2.

		When StructDepth is greater than 1 the resulting error label will be the "path" to the field within the top level struct.
		So if C in the example above did have an issue the error label would be "B.C".
	*/
	StructDepth uint8
	/*
		Value is the raw value being validated.
	*/
	Value interface{}
}

func New() ValidationParams {
	return ValidationParams{
		ArrayDepth:  0,
		Name:        "",
		Required:    false,
		StructDepth: 0,
	}
}
