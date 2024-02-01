package builtin

import (
	"errors"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	gbutil "github.com/Ghostbb-io/gb/util/gb_util"
	"strings"
)

// RuleSame implements `same` rule:
// Value should be the same as value of field.
//
// Format: same:field
type RuleSame struct{}

func init() {
	Register(RuleSame{})
}

func (r RuleSame) Name() string {
	return "same"
}

func (r RuleSame) Message() string {
	return "The {field} value `{value}` must be the same as field {field1} value `{value1}`"
}

func (r RuleSame) Run(in RunInput) error {
	var (
		ok    bool
		value = in.Value.String()
	)
	fieldName, fieldValue := gbutil.MapPossibleItemByKey(in.Data.Map(), in.RulePattern)
	if fieldValue != nil {
		if in.Option.CaseInsensitive {
			ok = strings.EqualFold(value, gbconv.String(fieldValue))
		} else {
			ok = strings.Compare(value, gbconv.String(fieldValue)) == 0
		}
	}
	if !ok {
		return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
			"{field1}": fieldName,
			"{value1}": gbconv.String(fieldValue),
		}))
	}
	return nil
}
