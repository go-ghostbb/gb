package builtin

import (
	"errors"
	gbregex "ghostbb.io/text/gb_regex"
)

// RulePostcode implements `postcode` rule:
// Postcode number.
//
// Format: postcode
type RulePostcode struct{}

func init() {
	Register(RulePostcode{})
}

func (r RulePostcode) Name() string {
	return "postcode"
}

func (r RulePostcode) Message() string {
	return "The {field} value `{value}` is not a valid postcode format"
}

func (r RulePostcode) Run(in RunInput) error {
	ok := gbregex.IsMatchString(
		`^\d{6}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
