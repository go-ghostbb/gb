package builtin

import (
	"errors"
	gbipv4 "ghostbb.io/gb/net/gb_ipv4"
	gbipv6 "ghostbb.io/gb/net/gb_ipv6"
)

// RuleIp implements `ip` rule:
// IPv4/IPv6.
//
// Format: ip
type RuleIp struct{}

func init() {
	Register(RuleIp{})
}

func (r RuleIp) Name() string {
	return "ip"
}

func (r RuleIp) Message() string {
	return "The {field} value `{value}` is not a valid IP address"
}

func (r RuleIp) Run(in RunInput) error {
	var (
		ok    bool
		value = in.Value.String()
	)
	if ok = gbipv4.Validate(value) || gbipv6.Validate(value); !ok {
		return errors.New(in.Message)
	}
	return nil
}
