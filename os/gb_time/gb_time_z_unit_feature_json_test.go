package gbtime_test

import (
	"ghostbb.io/gb/frame/g"
	"ghostbb.io/gb/internal/json"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_Json_Pointer(t *testing.T) {
	// Marshal
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime *gbtime.Time
		}
		b, err := json.Marshal(MyTime{
			MyTime: gbtime.NewFromStr("2006-01-02 15:04:05"),
		})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":"2006-01-02 15:04:05"}`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		b, err := json.Marshal(g.Map{
			"MyTime": gbtime.NewFromStr("2006-01-02 15:04:05"),
		})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":"2006-01-02 15:04:05"}`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		b, err := json.Marshal(g.Map{
			"MyTime": *gbtime.NewFromStr("2006-01-02 15:04:05"),
		})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime *gbtime.Time
		}
		b, err := json.Marshal(&MyTime{})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":null}`)
	})
	// Marshal nil with json omitempty
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime *gbtime.Time `json:"time,omitempty"`
		}
		b, err := json.Marshal(&MyTime{})
		t.AssertNil(err)
		t.Assert(b, `{}`)
	})
	// Unmarshal
	gbtest.C(t, func(t *gbtest.T) {
		var (
			myTime gbtime.Time
			err    = json.UnmarshalUseNumber([]byte(`"2006-01-02 15:04:05"`), &myTime)
		)
		t.AssertNil(err)
		t.Assert(myTime.String(), "2006-01-02 15:04:05")
	})
}

func Test_Json_Struct(t *testing.T) {
	// Marshal struct.
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime gbtime.Time
		}
		b, err := json.Marshal(MyTime{
			MyTime: *gbtime.NewFromStr("2006-01-02 15:04:05"),
		})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":"2006-01-02 15:04:05"}`)
	})
	// Marshal pointer.
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime gbtime.Time
		}
		b, err := json.Marshal(&MyTime{
			MyTime: *gbtime.NewFromStr("2006-01-02 15:04:05"),
		})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime gbtime.Time
		}
		b, err := json.Marshal(MyTime{})
		t.AssertNil(err)
		t.Assert(b, `{"MyTime":""}`)
	})
	// Marshal nil omitempty
	gbtest.C(t, func(t *gbtest.T) {
		type MyTime struct {
			MyTime gbtime.Time `json:"time,omitempty"`
		}
		b, err := json.Marshal(MyTime{})
		t.AssertNil(err)
		t.Assert(b, `{"time":""}`)
	})

}
