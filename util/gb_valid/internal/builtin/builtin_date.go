package builtin

import (
	"errors"
	gbregex "ghostbb.io/gb/text/gb_regex"
	"time"
)

// RuleDate implements `date` rule:
// Standard date, like: 2006-01-02, 20060102, 2006.01.02.
//
// Format: date
type RuleDate struct{}

func init() {
	Register(RuleDate{})
}

func (r RuleDate) Name() string {
	return "date"
}

func (r RuleDate) Message() string {
	return "The {field} value `{value}` is not a valid date"
}

func (r RuleDate) Run(in RunInput) error {
	type iTime interface {
		Date() (year int, month time.Month, day int)
		IsZero() bool
	}
	// support for time value, eg: gbtime.Time/*gbtime.Time, time.Time/*time.Time.
	if obj, ok := in.Value.Val().(iTime); ok {
		if obj.IsZero() {
			return errors.New(in.Message)
		}
	}
	if !gbregex.IsMatchString(
		`\d{4}[\.\-\_/]{0,1}\d{2}[\.\-\_/]{0,1}\d{2}`,
		in.Value.String(),
	) {
		return errors.New(in.Message)
	}
	return nil
}
