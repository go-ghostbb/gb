package builtin

import (
	"errors"
	gbipv4 "ghostbb.io/net/gb_ipv4"
)

// RuleIpv4 implements `ipv4` rule:
// IPv4.
//
// Format: ipv4
type RuleIpv4 struct{}

func init() {
	Register(RuleIpv4{})
}

func (r RuleIpv4) Name() string {
	return "ipv4"
}

func (r RuleIpv4) Message() string {
	return "The {field} value `{value}` is not a valid IPv4 address"
}

func (r RuleIpv4) Run(in RunInput) error {
	var (
		ok    bool
		value = in.Value.String()
	)
	if ok = gbipv4.Validate(value); !ok {
		return errors.New(in.Message)
	}
	return nil
}
