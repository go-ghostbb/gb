package builtin

import (
	"errors"
	"ghostbb.io/internal/json"
)

// RuleJson implements `json` rule:
// JSON.
//
// Format: json
type RuleJson struct{}

func init() {
	Register(RuleJson{})
}

func (r RuleJson) Name() string {
	return "json"
}

func (r RuleJson) Message() string {
	return "The {field} value `{value}` is not a valid JSON string"
}

func (r RuleJson) Run(in RunInput) error {
	if json.Valid(in.Value.Bytes()) {
		return nil
	}
	return errors.New(in.Message)
}
