package builtin

import (
	"errors"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	gbutil "github.com/Ghostbb-io/gb/util/gb_util"
	"strconv"
)

// RuleGT implements `gt` rule:
// Greater than `field`.
// It supports both integer and float.
//
// Format: gt:field
type RuleGT struct{}

func init() {
	Register(RuleGT{})
}

func (r RuleGT) Name() string {
	return "gt"
}

func (r RuleGT) Message() string {
	return "The {field} value `{value}` must be greater than field {field1} value `{value1}`"
}

func (r RuleGT) Run(in RunInput) error {
	var (
		fieldName, fieldValue = gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
		fieldValueN, err1     = strconv.ParseFloat(gbconv.String(fieldValue), 10)
		valueN, err2          = strconv.ParseFloat(in.Value.String(), 10)
	)

	if valueN <= fieldValueN || err1 != nil || err2 != nil {
		return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
			"{field1}": fieldName,
			"{value1}": gbconv.String(fieldValue),
		}))
	}
	return nil
}
