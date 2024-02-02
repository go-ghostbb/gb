package builtin

import (
	"errors"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"strings"
)

// RuleDifferent implements `different` rule:
// Value should be different from value of field.
//
// Format: different:field
type RuleDifferent struct{}

func init() {
	Register(RuleDifferent{})
}

func (r RuleDifferent) Name() string {
	return "different"
}

func (r RuleDifferent) Message() string {
	return "The {field} value `{value}` must be different from field {field1} value `{value1}`"
}

func (r RuleDifferent) Run(in RunInput) error {
	var (
		ok    = true
		value = in.Value.String()
	)
	fieldName, fieldValue := gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
	if fieldValue != nil {
		if in.Option.CaseInsensitive {
			ok = !strings.EqualFold(value, gbconv.String(fieldValue))
		} else {
			ok = strings.Compare(value, gbconv.String(fieldValue)) != 0
		}
	}
	if ok {
		return nil
	}
	return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
		"{field1}": fieldName,
		"{value1}": gbconv.String(fieldValue),
	}))
}
