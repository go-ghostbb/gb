package gbconv

import (
	"github.com/Ghostbb-io/gb/internal/utils"
	gbtime "github.com/Ghostbb-io/gb/os/gb_time"
	"time"
)

// Time converts `any` to time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	if t := GBTime(any, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !utils.IsNumeric(s) {
		d, _ := gbtime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}

// GBTime converts `any` to *gbtime.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using gtime.NewFromTimeStamp if `any` is numeric,
// or using gbtime.StrToTime if `any` is string.
func GBTime(any interface{}, format ...string) *gbtime.Time {
	if any == nil {
		return nil
	}
	if v, ok := any.(iGBTime); ok {
		return v.GBTime(format...)
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*gbtime.Time); ok {
			return v
		}
		if t, ok := any.(time.Time); ok {
			return gbtime.New(t)
		}
		if t, ok := any.(*time.Time); ok {
			return gbtime.New(t)
		}
	}
	s := String(any)
	if len(s) == 0 {
		return gbtime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := gbtime.StrToTimeFormat(s, item)
			if t != nil && err == nil {
				return t
			}
		}
		return nil
	}
	if utils.IsNumeric(s) {
		return gbtime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := gbtime.StrToTime(s)
		return t
	}
}
