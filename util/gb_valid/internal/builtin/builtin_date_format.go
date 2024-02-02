package builtin

import (
	"errors"
	gbtime "ghostbb.io/gb/os/gb_time"
	"time"
)

// RuleDateFormat implements `date-format` rule:
// Custom date format.
//
// Format: date-format:format
type RuleDateFormat struct{}

func init() {
	Register(RuleDateFormat{})
}

func (r RuleDateFormat) Name() string {
	return "date-format"
}

func (r RuleDateFormat) Message() string {
	return "The {field} value `{value}` does not match the format: {pattern}"
}

func (r RuleDateFormat) Run(in RunInput) error {
	type iTime interface {
		Date() (year int, month time.Month, day int)
		IsZero() bool
	}
	// support for time value, eg: gtime.Time/*gtime.Time, time.Time/*time.Time.
	if obj, ok := in.Value.Val().(iTime); ok {
		if obj.IsZero() {
			return errors.New(in.Message)
		}
		return nil
	}
	if _, err := gbtime.StrToTimeFormat(in.Value.String(), in.RulePattern); err != nil {
		return errors.New(in.Message)
	}
	return nil
}
