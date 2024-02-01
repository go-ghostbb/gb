package builtin

// RuleForeach implements `foreach` rule:
// It tells the next validation using current value as an array and validates each of its element.
//
// Format: foreach
type RuleForeach struct{}

func init() {
	Register(RuleForeach{})
}

func (r RuleForeach) Name() string {
	return "foreach"
}

func (r RuleForeach) Message() string {
	return ""
}

func (r RuleForeach) Run(in RunInput) error {
	return nil
}
