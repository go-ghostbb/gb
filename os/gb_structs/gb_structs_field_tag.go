package gbstructs

import (
	gbtag "github.com/Ghostbb-io/gb/util/gb_tag"
	"strings"
)

// TagJsonName returns the `json` tag name string of the field.
func (f *Field) TagJsonName() string {
	if jsonTag := f.Tag(gbtag.Json); jsonTag != "" {
		return strings.Split(jsonTag, ",")[0]
	}
	return ""
}

// TagDefault returns the most commonly used tag `default/d` value of the field.
func (f *Field) TagDefault() string {
	v := f.Tag(gbtag.Default)
	if v == "" {
		v = f.Tag(gbtag.DefaultShort)
	}
	return v
}

// TagParam returns the most commonly used tag `param/p` value of the field.
func (f *Field) TagParam() string {
	v := f.Tag(gbtag.Param)
	if v == "" {
		v = f.Tag(gbtag.ParamShort)
	}
	return v
}

// TagValid returns the most commonly used tag `valid/v` value of the field.
func (f *Field) TagValid() string {
	v := f.Tag(gbtag.Valid)
	if v == "" {
		v = f.Tag(gbtag.ValidShort)
	}
	return v
}

// TagDescription returns the most commonly used tag `description/des/dc` value of the field.
func (f *Field) TagDescription() string {
	v := f.Tag(gbtag.Description)
	if v == "" {
		v = f.Tag(gbtag.DescriptionShort)
	}
	if v == "" {
		v = f.Tag(gbtag.DescriptionShort2)
	}
	return v
}

// TagSummary returns the most commonly used tag `summary/sum/sm` value of the field.
func (f *Field) TagSummary() string {
	v := f.Tag(gbtag.Summary)
	if v == "" {
		v = f.Tag(gbtag.SummaryShort)
	}
	if v == "" {
		v = f.Tag(gbtag.SummaryShort2)
	}
	return v
}

// TagAdditional returns the most commonly used tag `additional/ad` value of the field.
func (f *Field) TagAdditional() string {
	v := f.Tag(gbtag.Additional)
	if v == "" {
		v = f.Tag(gbtag.AdditionalShort)
	}
	return v
}

// TagExample returns the most commonly used tag `example/eg` value of the field.
func (f *Field) TagExample() string {
	v := f.Tag(gbtag.Example)
	if v == "" {
		v = f.Tag(gbtag.ExampleShort)
	}
	return v
}

// TagIn returns the most commonly used tag `in` value of the field.
func (f *Field) TagIn() string {
	v := f.Tag(gbtag.In)
	return v
}
