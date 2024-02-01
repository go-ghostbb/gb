package builtin

import (
	"errors"
	gbstr "ghostbb.io/text/gb_str"
	gbconv "ghostbb.io/util/gb_conv"
	"strconv"
)

// RuleMinLength implements `min-length` rule:
// Length is equal or greater than :min.
// The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
//
// Format: min-length:min
type RuleMinLength struct{}

func init() {
	Register(RuleMinLength{})
}

func (r RuleMinLength) Name() string {
	return "min-length"
}

func (r RuleMinLength) Message() string {
	return "The {field} value `{value}` length must be equal or greater than {min}"
}

func (r RuleMinLength) Run(in RunInput) error {
	var (
		valueRunes = gbconv.Runes(in.Value.String())
		valueLen   = len(valueRunes)
	)
	min, err := strconv.Atoi(in.RulePattern)
	if valueLen < min || err != nil {
		return errors.New(gbstr.Replace(in.Message, "{min}", strconv.Itoa(min)))
	}
	return nil
}
