package builtin

import (
	"errors"
	gbstr "ghostbb.io/text/gb_str"
	"strconv"
	"strings"
)

// RuleBetween implements `between` rule:
// Range between :min and :max. It supports both integer and float.
//
// Format: between:min,max
type RuleBetween struct{}

func init() {
	Register(RuleBetween{})
}

func (r RuleBetween) Name() string {
	return "between"
}

func (r RuleBetween) Message() string {
	return "The {field} value `{value}` must be between {min} and {max}"
}

func (r RuleBetween) Run(in RunInput) error {
	var (
		array = strings.Split(in.RulePattern, ",")
		min   = float64(0)
		max   = float64(0)
	)
	if len(array) > 0 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(array[0]), 10); err == nil {
			min = v
		}
	}
	if len(array) > 1 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(array[1]), 10); err == nil {
			max = v
		}
	}
	valueF, err := strconv.ParseFloat(in.Value.String(), 10)
	if valueF < min || valueF > max || err != nil {
		return errors.New(gbstr.ReplaceByMap(in.Message, map[string]string{
			"{min}": strconv.FormatFloat(min, 'f', -1, 64),
			"{max}": strconv.FormatFloat(max, 'f', -1, 64),
		}))
	}
	return nil
}
