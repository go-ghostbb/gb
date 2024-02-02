package builtin

import (
	"errors"
	gbregex "ghostbb.io/gb/text/gb_regex"
)

// RulePassword2 implements `password2` rule:
// Universal password format rule2:
// Must meet password rule1, must contain lower and upper letters and numbers.
//
// Format: password2
type RulePassword2 struct{}

func init() {
	Register(RulePassword2{})
}

func (r RulePassword2) Name() string {
	return "password2"
}

func (r RulePassword2) Message() string {
	return "The {field} value `{value}` is not a valid password2 format"
}

func (r RulePassword2) Run(in RunInput) error {
	var value = in.Value.String()
	if gbregex.IsMatchString(`^[\w\S]{6,18}$`, value) &&
		gbregex.IsMatchString(`[a-z]+`, value) &&
		gbregex.IsMatchString(`[A-Z]+`, value) &&
		gbregex.IsMatchString(`\d+`, value) {
		return nil
	}
	return errors.New(in.Message)
}
