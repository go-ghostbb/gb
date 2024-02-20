package builtin

import (
	"errors"
	gbtime "ghostbb.io/gb/os/gb_time"
	"time"
)

// RuleDatetime implements `datetime` rule:
// Standard datetime, like: 2006-01-02 12:00:00.
//
// Format: datetime
type RuleDatetime struct{}

func init() {
	Register(RuleDatetime{})
}

func (r RuleDatetime) Name() string {
	return "datetime"
}

func (r RuleDatetime) Message() string {
	return "The {field} value `{value}` is not a valid datetime"
}

func (r RuleDatetime) Run(in RunInput) error {
	type iTime interface {
		Date() (year int, month time.Month, day int)
		IsZero() bool
	}
	// support for time value, eg: gbtime.Time/*gbtime.Time, time.Time/*time.Time.
	if obj, ok := in.Value.Val().(iTime); ok {
		if obj.IsZero() {
			return errors.New(in.Message)
		}
		return nil
	}
	if _, err := gbtime.StrToTimeFormat(in.Value.String(), `Y-m-d H:i:s`); err != nil {
		return errors.New(in.Message)
	}
	return nil
}
