package gblog

import (
	"bytes"
	"context"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

var (
	ctx = context.TODO()
)

func Test_Print(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Print(ctx, 1, 2, 3)
		l.Printf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "["), 0)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Debug(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Debug(ctx, 1, 2, 3)
		l.Debugf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), defaultLevelPrefixes[LEVEL_DEBU]), 2)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Info(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Info(ctx, 1, 2, 3)
		l.Infof(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), defaultLevelPrefixes[LEVEL_INFO]), 2)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Notice(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Notice(ctx, 1, 2, 3)
		l.Noticef(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), defaultLevelPrefixes[LEVEL_NOTI]), 2)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Warning(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Warning(ctx, 1, 2, 3)
		l.Warningf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), defaultLevelPrefixes[LEVEL_WARN]), 2)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Error(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Error(ctx, 1, 2, 3)
		l.Errorf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_LevelPrefix(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := New()
		t.Assert(l.GetLevelPrefix(LEVEL_DEBU), defaultLevelPrefixes[LEVEL_DEBU])
		t.Assert(l.GetLevelPrefix(LEVEL_INFO), defaultLevelPrefixes[LEVEL_INFO])
		t.Assert(l.GetLevelPrefix(LEVEL_NOTI), defaultLevelPrefixes[LEVEL_NOTI])
		t.Assert(l.GetLevelPrefix(LEVEL_WARN), defaultLevelPrefixes[LEVEL_WARN])
		t.Assert(l.GetLevelPrefix(LEVEL_ERRO), defaultLevelPrefixes[LEVEL_ERRO])
		t.Assert(l.GetLevelPrefix(LEVEL_CRIT), defaultLevelPrefixes[LEVEL_CRIT])
		l.SetLevelPrefix(LEVEL_DEBU, "debug")
		t.Assert(l.GetLevelPrefix(LEVEL_DEBU), "debug")
		l.SetLevelPrefixes(map[int]string{
			LEVEL_CRIT: "critical",
		})
		t.Assert(l.GetLevelPrefix(LEVEL_DEBU), "debug")
		t.Assert(l.GetLevelPrefix(LEVEL_INFO), defaultLevelPrefixes[LEVEL_INFO])
		t.Assert(l.GetLevelPrefix(LEVEL_NOTI), defaultLevelPrefixes[LEVEL_NOTI])
		t.Assert(l.GetLevelPrefix(LEVEL_WARN), defaultLevelPrefixes[LEVEL_WARN])
		t.Assert(l.GetLevelPrefix(LEVEL_ERRO), defaultLevelPrefixes[LEVEL_ERRO])
		t.Assert(l.GetLevelPrefix(LEVEL_CRIT), "critical")
	})
	gbtest.C(t, func(t *gbtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		l.SetWriter(buffer)
		l.Debug(ctx, "test1")
		t.Assert(gbstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_DEBU]), true)

		buffer.Reset()

		l.SetLevelPrefix(LEVEL_DEBU, "debug")
		l.Debug(ctx, "test2")
		t.Assert(gbstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_DEBU]), false)
		t.Assert(gbstr.Contains(buffer.String(), "debug"), true)

		buffer.Reset()
		l.SetLevelPrefixes(map[int]string{
			LEVEL_ERRO: "error",
		})
		l.Error(ctx, "test3")
		t.Assert(gbstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_ERRO]), false)
		t.Assert(gbstr.Contains(buffer.String(), "error"), true)
	})
}
