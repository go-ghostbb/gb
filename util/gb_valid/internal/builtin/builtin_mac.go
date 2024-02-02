package builtin

import (
	"errors"
	gbregex "ghostbb.io/gb/text/gb_regex"
)

// RuleMac implements `mac` rule:
// MAC.
//
// Format: mac
type RuleMac struct{}

func init() {
	Register(RuleMac{})
}

func (r RuleMac) Name() string {
	return "mac"
}

func (r RuleMac) Message() string {
	return "The {field} value `{value}` is not a valid MAC address"
}

func (r RuleMac) Run(in RunInput) error {
	ok := gbregex.IsMatchString(
		`^([0-9A-Fa-f]{2}[\-:]){5}[0-9A-Fa-f]{2}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
