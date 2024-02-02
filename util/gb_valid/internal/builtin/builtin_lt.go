package builtin

import (
	"errors"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"strconv"
)

// RuleLT implements `lt` rule:
// Lesser than `field`.
// It supports both integer and float.
//
// Format: lt:field
type RuleLT struct{}

func init() {
	Register(RuleLT{})
}

func (r RuleLT) Name() string {
	return "lt"
}

func (r RuleLT) Message() string {
	return "The {field} value `{value}` must be lesser than field {field1} value `{value1}`"
}

func (r RuleLT) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		fieldValueN, err1     = strconv.ParseFloat(gbconv.String(fieldValue), 10)
		valueN, err2          = strconv.ParseFloat(in.Value.String(), 10)
	)

	if valueN >= fieldValueN || err1 != nil || err2 != nil {
		return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
			"{field1}": fieldName,
			"{value1}": gbconv.String(fieldValue),
		}))
	}
	return nil
}
