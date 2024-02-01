package builtin

// RuleNotEq implements `not-eq` rule:
// Value should be different from value of field.
//
// Format: not-eq:field
type RuleNotEq struct{}

func init() {
	Register(RuleNotEq{})
}

func (r RuleNotEq) Name() string {
	return "not-eq"
}

func (r RuleNotEq) Message() string {
	return "The {field} value `{value}` must not be equal to field {field1} value `{value1}`"
}

func (r RuleNotEq) Run(in RunInput) error {
	return RuleDifferent{}.Run(in)
}
