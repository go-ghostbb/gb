package gblog

import (
	"bytes"
	gbtest "ghostbb.io/gb/test/gb_test"
	"strings"
	"testing"
)

func Test_SetConfigWithMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := New()
		m := map[string]interface{}{
			"path":     "/var/log",
			"level":    "all",
			"stdout":   false,
			"StStatus": 0,
		}
		err := l.SetConfigWithMap(m)
		t.AssertNil(err)
		t.Assert(l.config.Path, m["path"])
		t.Assert(l.config.Level, LEVEL_ALL)
		t.Assert(l.config.StdoutPrint, m["stdout"])
	})
}

func Test_SetConfigWithMap_LevelStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		m := map[string]interface{}{
			"level": "all",
		}
		err := l.SetConfigWithMap(m)
		t.AssertNil(err)

		l.SetWriter(buffer)

		l.Debug(ctx, "test")
		l.Warning(ctx, "test")
		t.Assert(strings.Contains(buffer.String(), "DEBU"), true)
		t.Assert(strings.Contains(buffer.String(), "WARN"), true)
	})

	gbtest.C(t, func(t *gbtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		m := map[string]interface{}{
			"level": "warn",
		}
		err := l.SetConfigWithMap(m)
		t.AssertNil(err)
		l.SetWriter(buffer)
		l.Debug(ctx, "test")
		l.Warning(ctx, "test")
		t.Assert(strings.Contains(buffer.String(), "DEBU"), false)
		t.Assert(strings.Contains(buffer.String(), "WARN"), true)
	})
}
