package builtin

import (
	"errors"
	gbipv6 "ghostbb.io/gb/net/gb_ipv6"
)

// RuleIpv6 implements `ipv6` rule:
// IPv6.
//
// Format: ipv6
type RuleIpv6 struct{}

func init() {
	Register(RuleIpv6{})
}

func (r RuleIpv6) Name() string {
	return "ipv6"
}

func (r RuleIpv6) Message() string {
	return "The {field} value `{value}` is not a valid IPv6 address"
}

func (r RuleIpv6) Run(in RunInput) error {
	var (
		ok    bool
		value = in.Value.String()
	)
	if ok = gbipv6.Validate(value); !ok {
		return errors.New(in.Message)
	}
	return nil
}
