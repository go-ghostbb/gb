package builtin

import (
	"errors"
	gbconv "ghostbb.io/util/gb_conv"
	gbutil "ghostbb.io/util/gb_util"
	"strings"
)

// RuleRequiredIf implements `required-if` rule:
// Required unless all given field and its value are equal.
//
// Format:  required-if:field,value,...
// Example: required-if: id,1,age,18
type RuleRequiredIf struct{}

func init() {
	Register(RuleRequiredIf{})
}

func (r RuleRequiredIf) Name() string {
	return "required-if"
}

func (r RuleRequiredIf) Message() string {
	return "The {field} field is required"
}

func (r RuleRequiredIf) Run(in RunInput) error {
	var (
		required   = false
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
	)
	// It supports multiple field and value pairs.
	if len(array)%2 == 0 {
		for i := 0; i < len(array); {
			tk := array[i]
			tv := array[i+1]
			_, foundValue = gbutil.MapPossibleItemByKey(in.Data.Map(), tk)
			if in.Option.CaseInsensitive {
				required = strings.EqualFold(tv, gbconv.String(foundValue))
			} else {
				required = strings.Compare(tv, gbconv.String(foundValue)) == 0
			}
			if required {
				break
			}
			i += 2
		}
	}

	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
