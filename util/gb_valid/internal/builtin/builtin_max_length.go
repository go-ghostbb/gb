package builtin

import (
	"errors"
	gbstr "ghostbb.io/text/gb_str"
	gbconv "ghostbb.io/util/gb_conv"
	"strconv"
)

// RuleMaxLength implements `max-length` rule:
// Length is equal or lesser than :max.
// The length is calculated using unicode string, which means one chinese character or letter both has the length of 1.
//
// Format: max-length:max
type RuleMaxLength struct{}

func init() {
	Register(RuleMaxLength{})
}

func (r RuleMaxLength) Name() string {
	return "max-length"
}

func (r RuleMaxLength) Message() string {
	return "The {field} value `{value}` length must be equal or lesser than {max}"
}

func (r RuleMaxLength) Run(in RunInput) error {
	var (
		valueRunes = gbconv.Runes(in.Value.String())
		valueLen   = len(valueRunes)
	)
	max, err := strconv.Atoi(in.RulePattern)
	if valueLen > max || err != nil {
		return errors.New(gbstr.Replace(in.Message, "{max}", strconv.Itoa(max)))
	}
	return nil
}
