package builtin

import (
	"errors"
	"strconv"
)

// RuleInteger implements `integer` rule:
// Integer.
//
// Format: integer
type RuleInteger struct{}

func init() {
	Register(RuleInteger{})
}

func (r RuleInteger) Name() string {
	return "integer"
}

func (r RuleInteger) Message() string {
	return "The {field} value `{value}` is not an integer"
}

func (r RuleInteger) Run(in RunInput) error {
	if _, err := strconv.Atoi(in.Value.String()); err == nil {
		return nil
	}
	return errors.New(in.Message)
}
