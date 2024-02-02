package builtin

import (
	"errors"
	"fmt"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/json"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbtag "ghostbb.io/gb/util/gb_tag"
)

// RuleEnums implements `enums` rule:
// Value should be in enums of its constant type.
//
// Format: enums
type RuleEnums struct{}

func init() {
	Register(RuleEnums{})
}

func (r RuleEnums) Name() string {
	return "enums"
}

func (r RuleEnums) Message() string {
	return "The {field} value `{value}` should be in enums of: {enums}"
}

func (r RuleEnums) Run(in RunInput) error {
	if in.ValueType == nil {
		return gberror.NewCode(
			gbcode.CodeInvalidParameter,
			`value type cannot be empty to use validation rule "enums"`,
		)
	}
	var (
		pkgPath  = in.ValueType.PkgPath()
		typeName = in.ValueType.Name()
	)
	if pkgPath == "" {
		return gberror.NewCodef(
			gbcode.CodeInvalidOperation,
			`no pkg path found for type "%s"`,
			in.ValueType.String(),
		)
	}
	var (
		typeId   = fmt.Sprintf(`%s.%s`, pkgPath, typeName)
		tagEnums = gbtag.GetEnumsByType(typeId)
	)
	if tagEnums == "" {
		return gberror.NewCodef(
			gbcode.CodeInvalidOperation,
			`no enums found for type "%s"`,
			typeId,
		)
	}
	var enumsValues = make([]interface{}, 0)
	if err := json.Unmarshal([]byte(tagEnums), &enumsValues); err != nil {
		return err
	}
	if !gbstr.InArray(gbconv.Strings(enumsValues), in.Value.String()) {
		return errors.New(gbstr.Replace(
			in.Message, `{enums}`, tagEnums,
		))
	}
	return nil
}
