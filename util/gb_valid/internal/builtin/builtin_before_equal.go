package builtin

import (
	"errors"
	gbstr "ghostbb.io/text/gb_str"
	gbconv "ghostbb.io/util/gb_conv"
	gbutil "ghostbb.io/util/gb_util"
)

// RuleBeforeEqual implements `before-equal` rule:
// The datetime value should be after or equal to the value of field `field`.
//
// Format: before-equal:field
type RuleBeforeEqual struct{}

func init() {
	Register(RuleBeforeEqual{})
}

func (r RuleBeforeEqual) Name() string {
	return "before-equal"
}

func (r RuleBeforeEqual) Message() string {
	return "The {field} value `{value}` must be before or equal to field {pattern}"
}

func (r RuleBeforeEqual) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		valueDatetime         = in.Value.Time()
		fieldDatetime         = gbconv.Time(fieldValue)
	)
	if valueDatetime.Before(fieldDatetime) || valueDatetime.Equal(fieldDatetime) {
		return nil
	}
	return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
		"{field1}": fieldName,
		"{value1}": gbconv.String(fieldValue),
	}))
}
