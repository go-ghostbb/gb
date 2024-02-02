package builtin

import (
	"errors"
	gbregex "ghostbb.io/gb/text/gb_regex"
)

// RuleTelephone implements `telephone` rule:
// "XXXX-XXXXXXX"
// "XXXX-XXXXXXXX"
// "XXX-XXXXXXX"
// "XXX-XXXXXXXX"
// "XXXXXXX"
// "XXXXXXXX"
//
// Format: telephone
type RuleTelephone struct{}

func init() {
	Register(RuleTelephone{})
}

func (r RuleTelephone) Name() string {
	return "telephone"
}

func (r RuleTelephone) Message() string {
	return "The {field} value `{value}` is not a valid telephone number"
}

func (r RuleTelephone) Run(in RunInput) error {
	ok := gbregex.IsMatchString(
		`^((\d{3,4})|\d{3,4}-)?\d{7,8}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
