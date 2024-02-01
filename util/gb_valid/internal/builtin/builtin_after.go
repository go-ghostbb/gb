package builtin

import (
	"errors"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	gbutil "github.com/Ghostbb-io/gb/util/gb_util"
)

// RuleAfter implements `after` rule:
// The datetime value should be after the value of field `field`.
//
// Format: after:field
type RuleAfter struct{}

func init() {
	Register(RuleAfter{})
}

func (r RuleAfter) Name() string {
	return "after"
}

func (r RuleAfter) Message() string {
	return "The {field} value `{value}` must be after field {field1} value `{value1}`"
}

func (r RuleAfter) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		valueDatetime         = in.Value.Time()
		fieldDatetime         = gbconv.Time(fieldValue)
	)
	if valueDatetime.After(fieldDatetime) {
		return nil
	}
	return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
		"{field1}": fieldName,
		"{value1}": gbconv.String(fieldValue),
	}))
}
