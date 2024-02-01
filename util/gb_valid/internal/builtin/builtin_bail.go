package builtin

// RuleBail implements `bail` rule:
// Stop validating when this field's validation failed.
//
// Format: bail
type RuleBail struct{}

func init() {
	Register(RuleBail{})
}

func (r RuleBail) Name() string {
	return "bail"
}

func (r RuleBail) Message() string {
	return ""
}

func (r RuleBail) Run(in RunInput) error {
	return nil
}
