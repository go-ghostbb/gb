package builtin

import (
	"errors"
	gbregex "github.com/Ghostbb-io/gb/text/gb_regex"
)

// RuleQQ implements `qq` rule:
// Tencent QQ number.
//
// Format: qq
type RuleQQ struct{}

func init() {
	Register(RuleQQ{})
}

func (r RuleQQ) Name() string {
	return "qq"
}

func (r RuleQQ) Message() string {
	return "The {field} value `{value}` is not a valid QQ number"
}

func (r RuleQQ) Run(in RunInput) error {
	ok := gbregex.IsMatchString(
		`^[1-9][0-9]{4,}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
