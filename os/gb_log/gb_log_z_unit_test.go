package gblog_test

import (
	"bytes"
	"context"
	"ghostbb.io/gb/frame/g"
	gbfile "ghostbb.io/gb/os/gb_file"
	gblog "ghostbb.io/gb/os/gb_log"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCase(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)

	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNE(gblog.Instance(), nil)
	})
}

func TestDefaultLogger(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)

	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNE(defaultLog, nil)
		log := gblog.New()
		gblog.SetDefaultLogger(log)
		t.AssertEQ(gblog.DefaultLogger(), defaultLog)
		t.AssertEQ(gblog.Expose(), defaultLog)
	})
}

func TestAPI(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gblog.Print(ctx, "Print")
		gblog.Printf(ctx, "%s", "Printf")
		gblog.Info(ctx, "Info")
		gblog.Infof(ctx, "%s", "Infof")
		gblog.Debug(ctx, "Debug")
		gblog.Debugf(ctx, "%s", "Debugf")
		gblog.Notice(ctx, "Notice")
		gblog.Noticef(ctx, "%s", "Noticef")
		gblog.Warning(ctx, "Warning")
		gblog.Warningf(ctx, "%s", "Warningf")
		gblog.Error(ctx, "Error")
		gblog.Errorf(ctx, "%s", "Errorf")
		gblog.Critical(ctx, "Critical")
		gblog.Criticalf(ctx, "%s", "Criticalf")
	})
}

func TestChaining(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)

	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNE(gblog.Cat("module"), nil)
		t.AssertNE(gblog.File("test.log"), nil)
		t.AssertNE(gblog.Level(gblog.LEVEL_ALL), nil)
		t.AssertNE(gblog.LevelStr("all"), nil)
		t.AssertNE(gblog.Skip(1), nil)
		t.AssertNE(gblog.Stack(false), nil)
		t.AssertNE(gblog.StackWithFilter("none"), nil)
		t.AssertNE(gblog.Stdout(false), nil)
		t.AssertNE(gblog.Header(false), nil)
		t.AssertNE(gblog.Line(false), nil)
		t.AssertNE(gblog.Async(false), nil)
	})
}

func Test_SetFile(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetFile("test.log")
	})
}

func Test_SetTimeFormat(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)

		l.SetTimeFormat("2006-01-02T15:04:05.000Z07:00")
		l.Debug(ctx, "test")

		t.AssertGE(len(strings.Split(w.String(), "[DEBU]")), 1)
		datetime := strings.Trim(strings.Split(w.String(), "[DEBU]")[0], " ")

		_, err := time.Parse("2006-01-02T15:04:05.000Z07:00", datetime)
		t.AssertNil(err)
		_, err = time.Parse("2006-01-02 15:04:05.000", datetime)
		t.AssertNE(err, nil)
		_, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", datetime)
		t.AssertNE(err, nil)
	})
}

func Test_SetLevel(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetLevel(gblog.LEVEL_ALL)
		t.Assert(gblog.GetLevel()&gblog.LEVEL_ALL, gblog.LEVEL_ALL)
	})
}

func Test_SetAsync(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetAsync(false)
	})
}

func Test_SetStdoutPrint(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetStdoutPrint(false)
	})
}

func Test_SetHeaderPrint(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetHeaderPrint(false)
	})
}

func Test_SetPrefix(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetPrefix("log_prefix")
	})
}

func Test_SetConfigWithMap(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gblog.SetConfigWithMap(map[string]interface{}{
			"level": "all",
		}), nil)
	})
}

func Test_SetPath(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gblog.SetPath("/var/log"), nil)
		t.Assert(gblog.GetPath(), "/var/log")
	})
}

func Test_SetWriter(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetWriter(os.Stdout)
		t.Assert(gblog.GetWriter(), os.Stdout)
	})
}

func Test_SetFlags(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetFlags(gblog.F_ASYNC)
		t.Assert(gblog.GetFlags(), gblog.F_ASYNC)
	})
}

func Test_SetCtxKeys(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetCtxKeys("SpanId", "TraceId")
		t.Assert(gblog.GetCtxKeys(), []string{"SpanId", "TraceId"})
	})
}

func Test_PrintStack(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.PrintStack(ctx, 1)
	})
}

func Test_SetStack(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetStack(true)
		t.Assert(gblog.GetStack(1), "")
	})
}

func Test_SetLevelStr(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gblog.SetLevelStr("all"), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		t.AssertNE(l.SetLevelStr("test"), nil)
	})
}

func Test_SetLevelPrefix(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetLevelPrefix(gblog.LEVEL_ALL, "LevelPrefix")
		t.Assert(gblog.GetLevelPrefix(gblog.LEVEL_ALL), "LevelPrefix")
	})
}

func Test_SetLevelPrefixes(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetLevelPrefixes(map[int]string{
			gblog.LEVEL_ALL: "ALL_Prefix",
		})
	})
}

func Test_SetHandlers(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetHandlers(func(ctx context.Context, in *gblog.HandlerInput) {
		})
	})
}

func Test_SetWriterColorEnable(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		gblog.SetWriterColorEnable(true)
	})
}

func Test_Instance(t *testing.T) {
	defaultLog := gblog.DefaultLogger().Clone()
	defer gblog.SetDefaultLogger(defaultLog)
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNE(gblog.Instance("gf"), nil)
	})
}

func Test_GetConfig(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		config := gblog.DefaultLogger().GetConfig()
		t.Assert(config.Path, "")
		t.Assert(config.StdoutPrint, true)
	})
}

func Test_Write(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		len, err := l.Write([]byte("Ghostbb"))
		t.AssertNil(err)
		t.Assert(len, 7)
	})
}

func Test_Chaining_To(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.DefaultLogger().Clone()
		logTo := l.To(os.Stdout)
		t.AssertNE(logTo, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logTo := l.To(os.Stdout)
		t.AssertNE(logTo, nil)
	})
}

func Test_Chaining_Path(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.DefaultLogger().Clone()
		logPath := l.Path("./")
		t.AssertNE(logPath, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logPath := l.Path("./")
		t.AssertNE(logPath, nil)
	})
}

func Test_Chaining_Cat(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logCat := l.Cat(".gf")
		t.AssertNE(logCat, nil)
	})
}

func Test_Chaining_Level(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logLevel := l.Level(gblog.LEVEL_ALL)
		t.AssertNE(logLevel, nil)
	})
}

func Test_Chaining_LevelStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logLevelStr := l.LevelStr("all")
		t.AssertNE(logLevelStr, nil)
	})
}

func Test_Chaining_Skip(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logSkip := l.Skip(1)
		t.AssertNE(logSkip, nil)
	})
}

func Test_Chaining_Stack(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logStack := l.Stack(true)
		t.AssertNE(logStack, nil)
	})
}

func Test_Chaining_StackWithFilter(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logStackWithFilter := l.StackWithFilter("gbtest")
		t.AssertNE(logStackWithFilter, nil)
	})
}

func Test_Chaining_Stdout(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logStdout := l.Stdout(true)
		t.AssertNE(logStdout, nil)
	})
}

func Test_Chaining_Header(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logHeader := l.Header(true)
		t.AssertNE(logHeader, nil)
	})
}

func Test_Chaining_Line(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logLine := l.Line(true)
		t.AssertNE(logLine, nil)
	})
}

func Test_Chaining_Async(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		logAsync := l.Async(true)
		t.AssertNE(logAsync, nil)
	})
}

func Test_Config_SetDebug(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		l.SetDebug(false)
	})
}

func Test_Config_AppendCtxKeys(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		l.AppendCtxKeys("Trace-Id", "Span-Id", "Test")
		l.AppendCtxKeys("Trace-Id-New", "Span-Id-New", "Test")
	})
}

func Test_Config_SetPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		t.AssertNE(l.SetPath(""), nil)
	})
}

func Test_Config_SetStdoutColorDisabled(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		l.SetStdoutColorDisabled(false)
	})
}

func Test_Ctx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gbstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 1)
	})
}

func Test_Ctx_Config(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		m := map[string]interface{}{
			"CtxKeys": g.SliceStr{"Trace-Id", "Span-Id", "Test"},
		}
		var nilMap map[string]interface{}

		err := l.SetConfigWithMap(m)
		t.AssertNil(err)
		err = l.SetConfigWithMap(nilMap)
		t.AssertNE(err, nil)

		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gbstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 1)
	})
}

func Test_Concurrent(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c := 1000
		l := gblog.New()
		s := "@1234567890#"
		f := "test.log"
		p := gbfile.Temp(gbtime.TimestampNanoStr())
		t.Assert(l.SetPath(p), nil)
		defer gbfile.Remove(p)
		wg := sync.WaitGroup{}
		ch := make(chan struct{})
		for i := 0; i < c; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ch
				l.File(f).Stdout(false).Print(ctx, s)
			}()
		}
		close(ch)
		wg.Wait()
		content := gbfile.GetContents(gbfile.Join(p, f))
		t.Assert(gbstr.Count(content, s), c)
	})
}
