package builtin

import (
	"errors"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	gbutil "github.com/Ghostbb-io/gb/util/gb_util"
)

// RuleBefore implements `before` rule:
// The datetime value should be after the value of field `field`.
//
// Format: before:field
type RuleBefore struct{}

func init() {
	Register(RuleBefore{})
}

func (r RuleBefore) Name() string {
	return "before"
}

func (r RuleBefore) Message() string {
	return "The {field} value `{value}` must be before field {field1} value `{value1}`"
}

func (r RuleBefore) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		valueDatetime         = in.Value.Time()
		fieldDatetime         = gbconv.Time(fieldValue)
	)
	if valueDatetime.Before(fieldDatetime) {
		return nil
	}
	return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
		"{field1}": fieldName,
		"{value1}": gbconv.String(fieldValue),
	}))
}
