package builtin

import (
	"errors"
	gbregex "ghostbb.io/gb/text/gb_regex"
)

// RulePhoneLoose implements `phone-loose` rule:
// Loose mobile phone number verification(宽松的手机号验证)
// As long as the 11 digits numbers beginning with
// 13, 14, 15, 16, 17, 18, 19 can pass the verification
// (只要满足 13、14、15、16、17、18、19开头的11位数字都可以通过验证).
//
// Format: phone-loose
type RulePhoneLoose struct{}

func init() {
	Register(RulePhoneLoose{})
}

func (r RulePhoneLoose) Name() string {
	return "phone-loose"
}

func (r RulePhoneLoose) Message() string {
	return "The {field} value `{value}` is not a valid phone number"
}

func (r RulePhoneLoose) Run(in RunInput) error {
	ok := gbregex.IsMatchString(
		`^1(3|4|5|6|7|8|9)\d{9}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
