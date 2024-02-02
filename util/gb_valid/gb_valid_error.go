package gbvalid

import (
	gbset "ghostbb.io/gb/container/gb_set"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbstr "ghostbb.io/gb/text/gb_str"
	"strings"
)

// Error is the validation error for validation result.
type Error interface {
	Code() gbcode.Code
	Current() error
	Error() string
	FirstItem() (key string, messages map[string]error)
	FirstRule() (rule string, err error)
	FirstError() (err error)
	Items() (items []map[string]map[string]error)
	Map() map[string]error
	Maps() map[string]map[string]error
	String() string
	Strings() (errs []string)
}

// validationError is the validation error for validation result.
type validationError struct {
	code      gbcode.Code                 // Error code.
	rules     []fieldRule                 // Rules by sequence, which is used for keeping error sequence only.
	errors    map[string]map[string]error // Error map:map[field]map[rule]message
	firstKey  string                      // The first error rule key(empty in default).
	firstItem map[string]error            // The first error rule value(nil in default).
}

// newValidationError creates and returns a validation error.
func newValidationError(code gbcode.Code, rules []fieldRule, fieldRuleErrorMap map[string]map[string]error) *validationError {
	for field, ruleErrorMap := range fieldRuleErrorMap {
		for rule, err := range ruleErrorMap {
			if !gberror.HasStack(err) {
				ruleErrorMap[rule] = gberror.NewWithOption(gberror.Option{
					Stack: false,
					Text:  gbstr.Trim(err.Error()),
					Code:  code,
				})
			}
		}
		fieldRuleErrorMap[field] = ruleErrorMap
	}
	// Filter repeated sequence rules.
	var ruleNameSet = gbset.NewStrSet()
	for i := 0; i < len(rules); {
		if !ruleNameSet.AddIfNotExist(rules[i].Name) {
			// Delete repeated rule.
			rules = append(rules[:i], rules[i+1:]...)
			continue
		}
		i++
	}
	return &validationError{
		code:   code,
		rules:  rules,
		errors: fieldRuleErrorMap,
	}
}

// newValidationErrorByStr creates and returns a validation error by string.
func newValidationErrorByStr(key string, err error) *validationError {
	return newValidationError(
		gbcode.CodeInternalError,
		nil,
		map[string]map[string]error{
			internalErrorMapKey: {
				key: err,
			},
		},
	)
}

// Code returns the error code of current validation error.
func (e *validationError) Code() gbcode.Code {
	if e == nil {
		return gbcode.CodeNil
	}
	return e.code
}

// Map returns the first error message as map.
func (e *validationError) Map() map[string]error {
	if e == nil {
		return map[string]error{}
	}
	_, m := e.FirstItem()
	return m
}

// Maps returns all error messages as map.
func (e *validationError) Maps() map[string]map[string]error {
	if e == nil {
		return nil
	}
	return e.errors
}

// Items retrieves and returns error items array in sequence if possible,
// or else it returns error items with no sequence .
func (e *validationError) Items() (items []map[string]map[string]error) {
	if e == nil {
		return []map[string]map[string]error{}
	}
	items = make([]map[string]map[string]error, 0)
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				items = append(items, map[string]map[string]error{
					v.Name: errorItemMap,
				})
			}
		}
		return items
	}
	// No sequence.
	for name, errorRuleMap := range e.errors {
		items = append(items, map[string]map[string]error{
			name: errorRuleMap,
		})
	}
	return
}

// FirstItem returns the field name and error messages for the first validation rule error.
func (e *validationError) FirstItem() (key string, messages map[string]error) {
	if e == nil {
		return "", map[string]error{}
	}
	if e.firstItem != nil {
		return e.firstKey, e.firstItem
	}
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				e.firstKey = v.Name
				e.firstItem = errorItemMap
				return v.Name, errorItemMap
			}
		}
	}
	// No sequence.
	for k, m := range e.errors {
		e.firstKey = k
		e.firstItem = m
		return k, m
	}
	return "", nil
}

// FirstRule returns the first error rule and message string.
func (e *validationError) FirstRule() (rule string, err error) {
	if e == nil {
		return "", nil
	}
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				for _, ruleItem := range strings.Split(v.Rule, "|") {
					array := strings.Split(ruleItem, ":")
					ruleItem = strings.TrimSpace(array[0])
					if err, ok = errorItemMap[ruleItem]; ok {
						return ruleItem, err
					}
				}
			}
		}
	}
	// No sequence.
	for _, errorItemMap := range e.errors {
		for k, v := range errorItemMap {
			return k, v
		}
	}
	return "", nil
}

// FirstError returns the first error message as string.
// Note that the returned message might be different if it has no sequence.
func (e *validationError) FirstError() (err error) {
	if e == nil {
		return nil
	}
	_, err = e.FirstRule()
	return
}

// Current is alis of FirstError, which implements interface gberror.iCurrent.
func (e *validationError) Current() error {
	return e.FirstError()
}

// String returns all error messages as string, multiple error messages joined using char ';'.
func (e *validationError) String() string {
	if e == nil {
		return ""
	}
	return strings.Join(e.Strings(), "; ")
}

// Error implements interface of error.Error.
func (e *validationError) Error() string {
	if e == nil {
		return ""
	}
	return e.String()
}

// Strings returns all error messages as string array.
func (e *validationError) Strings() (errs []string) {
	if e == nil {
		return []string{}
	}
	errs = make([]string, 0)
	// By sequence.
	if len(e.rules) > 0 {
		for _, v := range e.rules {
			if errorItemMap, ok := e.errors[v.Name]; ok {
				// validation error checks.
				for _, ruleItem := range strings.Split(v.Rule, "|") {
					ruleItem = strings.TrimSpace(strings.Split(ruleItem, ":")[0])
					if err, ok := errorItemMap[ruleItem]; ok {
						errs = append(errs, err.Error())
					}
				}
				// internal error checks.
				for k := range internalErrKeyMap {
					if err, ok := errorItemMap[k]; ok {
						errs = append(errs, err.Error())
					}
				}
			}
		}
		return errs
	}
	// No sequence.
	for _, errorItemMap := range e.errors {
		for _, err := range errorItemMap {
			errs = append(errs, err.Error())
		}
	}
	return
}
