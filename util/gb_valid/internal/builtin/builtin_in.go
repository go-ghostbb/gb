package builtin

import (
	"errors"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	"strings"
)

// RuleIn implements `in` rule:
// Value should be in: value1,value2,...
//
// Format: in:value1,value2,...
type RuleIn struct{}

func init() {
	Register(RuleIn{})
}

func (r RuleIn) Name() string {
	return "in"
}

func (r RuleIn) Message() string {
	return "The {field} value `{value}` is not in acceptable range: {pattern}"
}

func (r RuleIn) Run(in RunInput) error {
	var ok bool
	for _, rulePattern := range gbstr.SplitAndTrim(in.RulePattern, ",") {
		if in.Option.CaseInsensitive {
			ok = strings.EqualFold(in.Value.String(), strings.TrimSpace(rulePattern))
		} else {
			ok = strings.Compare(in.Value.String(), strings.TrimSpace(rulePattern)) == 0
		}
		if ok {
			return nil
		}
	}
	return errors.New(in.Message)
}
