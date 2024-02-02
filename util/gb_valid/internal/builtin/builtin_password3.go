package builtin

import (
	"errors"
	gbregex "ghostbb.io/gb/text/gb_regex"
)

// RulePassword3 implements `password3` rule:
// Universal password format rule3:
// Must meet password rule1, must contain lower and upper letters, numbers and special chars.
//
// Format: password3
type RulePassword3 struct{}

func init() {
	Register(RulePassword3{})
}

func (r RulePassword3) Name() string {
	return "password3"
}

func (r RulePassword3) Message() string {
	return "The {field} value `{value}` is not a valid password3 format"
}

func (r RulePassword3) Run(in RunInput) error {
	var value = in.Value.String()
	if gbregex.IsMatchString(`^[\w\S]{6,18}$`, value) &&
		gbregex.IsMatchString(`[a-z]+`, value) &&
		gbregex.IsMatchString(`[A-Z]+`, value) &&
		gbregex.IsMatchString(`\d+`, value) &&
		gbregex.IsMatchString(`[^a-zA-Z0-9]+`, value) {
		return nil
	}
	return errors.New(in.Message)
}
