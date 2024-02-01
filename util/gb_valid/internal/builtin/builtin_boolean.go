package builtin

import (
	"errors"
	"strings"
)

// RuleBoolean implements `boolean` rule:
// Boolean(1,true,on,yes:true | 0,false,off,no,"":false)
//
// Format: boolean
type RuleBoolean struct{}

// boolMap defines the boolean values.
var boolMap = map[string]struct{}{
	"1":     {},
	"true":  {},
	"on":    {},
	"yes":   {},
	"":      {},
	"0":     {},
	"false": {},
	"off":   {},
	"no":    {},
}

func init() {
	Register(RuleBoolean{})
}

func (r RuleBoolean) Name() string {
	return "boolean"
}

func (r RuleBoolean) Message() string {
	return "The {field} value `{value}` field must be true or false"
}

func (r RuleBoolean) Run(in RunInput) error {
	if _, ok := boolMap[strings.ToLower(in.Value.String())]; ok {
		return nil
	}
	return errors.New(in.Message)
}
