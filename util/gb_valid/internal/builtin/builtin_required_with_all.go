package builtin

import (
	"errors"
	"github.com/Ghostbb-io/gb/internal/empty"
	gbutil "github.com/Ghostbb-io/gb/util/gb_util"
	"strings"
)

// RuleRequiredWithAll implements `required-with-all` rule:
// Required if all given fields are not empty.
//
// Format:  required-with-all:field1,field2,...
// Example: required-with-all:id,name
type RuleRequiredWithAll struct{}

func init() {
	Register(RuleRequiredWithAll{})
}

func (r RuleRequiredWithAll) Name() string {
	return "required-with-all"
}

func (r RuleRequiredWithAll) Message() string {
	return "The {field} field is required"
}

func (r RuleRequiredWithAll) Run(in RunInput) error {
	var (
		required   = true
		array      = strings.Split(in.RulePattern, ",")
		foundValue interface{}
	)
	for i := 0; i < len(array); i++ {
		_, foundValue = gbutil.MapPossibleItemByKey(in.Data.Map(), array[i])
		if empty.IsEmpty(foundValue) {
			required = false
			break
		}
	}

	if required && isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}
