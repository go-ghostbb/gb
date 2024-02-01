package builtin

import (
	"errors"
	gbstr "ghostbb.io/text/gb_str"
	gbconv "ghostbb.io/util/gb_conv"
	gbutil "ghostbb.io/util/gb_util"
)

// RuleAfterEqual implements `after-equal` rule:
// The datetime value should be after or equal to the value of field `field`.
//
// Format: after-equal:field
type RuleAfterEqual struct{}

func init() {
	Register(RuleAfterEqual{})
}

func (r RuleAfterEqual) Name() string {
	return "after-equal"
}

func (r RuleAfterEqual) Message() string {
	return "The {field} value `{value}` must be after or equal to field {field1} value `{value1}`"
}

func (r RuleAfterEqual) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		valueDatetime         = in.Value.Time()
		fieldDatetime         = gbconv.Time(fieldValue)
	)
	if valueDatetime.After(fieldDatetime) || valueDatetime.Equal(fieldDatetime) {
		return nil
	}
	return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
		"{field1}": fieldName,
		"{value1}": gbconv.String(fieldValue),
	}))
}
