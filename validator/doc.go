/*
	The validator package contains the validators use by the validation package to perform validation.

	The Validator interface defined in this package require that a function ReadOptionsFromTagItems([]string) error be implemented.

	The proper syntax for struct field tag validator instructions is as follows:

		`validate:"validatorname,required,validatorspecificvalues"`

	So for instance using the int validator might have a tag lie this:

		`validate:"int,min=3,max=15"`

	This could be read by the validation package to produce an array of validation arguments from a struct field tag.
	The format is a comma sperated list of validation parameters where:

		- The first parameter is always the validator name as registered in the validators map in the validation package.
			- There is a special case when validating an embedded struct, you need the tag data, but for the validator name you need to add "struct" like below:
				- `validate:"struct"`
			- You can still have the required parameter in the tag data also, and if the underlying field is a pointer then the normal required rules for a pointer apply.
		- The second parameter is the validator parameters, if you are using the required parameter for any validator, it bus the the second parameter in the tag data to be registered properly.
		- After than, any additional validator parameters that you may need

	The validator parameters are the read by the above mentioned ReadOptionsFromTagItems function implemented by the validator matched by the validator name is the tag data.

	For examples of how a validator is implemented take a look at the various validators implemented in this package.
*/
package validator
