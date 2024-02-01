package gbdb

import (
	"context"
	"errors"
	"fmt"
	gblog "github.com/Ghostbb-io/gb/os/gb_log"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type ILogger interface {
	SetSlowThreshold(t time.Duration)
	SetIgnoreRecordNotFoundError(b bool)
	SetLogCat(s string)
	SetLogStdout(b bool)
	LogMode(gormlogger.LogLevel) gormlogger.Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

func NewLogger(group string, logger *gblog.Logger) ILogger {
	return &Logger{
		logger: logger,
		group:  group,
	}
}

type Logger struct {
	logger *gblog.Logger
	group  string
	LogConfig
}

type LogConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  gormlogger.LogLevel
	LogCat                    string
	LogStdout                 bool
}

func (l *Logger) SetSlowThreshold(t time.Duration) {
	l.SlowThreshold = t
}

func (l *Logger) SetIgnoreRecordNotFoundError(b bool) {
	l.IgnoreRecordNotFoundError = b
}

func (l *Logger) SetLogCat(s string) {
	l.LogCat = s
}

func (l *Logger) SetLogStdout(b bool) {
	if globalLogger.GetConfig().StdoutPrint {
		l.LogStdout = b
	} else {
		l.LogStdout = false
	}
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.LogLevel = level
	return l
}

func (l *Logger) Info(ctx context.Context, str string, args ...interface{}) {
	logger := l.setLogger()
	logger.Infof(ctx, l.group+" "+str, args...)
}

func (l *Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	logger := l.setLogger()
	logger.Warningf(ctx, l.group+" "+str, args...)
}

func (l *Logger) Error(ctx context.Context, str string, args ...interface{}) {
	logger := l.setLogger()
	logger.Errorf(ctx, l.group+" "+str, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	var (
		elapsed   = time.Since(begin)
		logger    = l.setLogger()
		track     = l.pathTrack(ctx)
		sql, rows = fc()
		msg       = fmt.Sprintf(`%s 【%s】｜%10v｜%10s`,
			track,
			l.group,
			elapsed,
			fmt.Sprintf("rows:%d", rows),
		)
		elapsedColor = l.ElapsedColor(elapsed)
		resetColor   = reset
	)
	sql = gbstr.Replace(sql, `/`, "", -1)
	msg += "\n    " + sql + "\n"

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		logger.Error(ctx, msg, map[string]interface{}{"error": err.Error()})
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		logger.Warning(ctx, msg)
	case l.LogLevel >= gormlogger.Info:
		logger.Info(ctx, msg)
	}

	fmt.Println(fmt.Sprintf(`%s [ORM] %s ｜%s %15v %s｜%s %-8s %s｜%-10s｜ %s`,
		time.Now().Format("2006/01/02 15:04:05"),
		track,
		elapsedColor, elapsed, resetColor,
		magenta, fmt.Sprintf("rows:%d", rows), reset,
		l.group,
		sql,
	))
}

func (l *Logger) ElapsedColor(elapsed time.Duration) string {
	if elapsed > l.SlowThreshold {
		return red
	}
	return green
}

const (
	gormPackage = "gorm.io"
	gbdbPackage = "database/gb_db"
	testFile    = "_test.go"
)

func (l *Logger) pathTrack(ctx context.Context) string {
	// Caller lookup
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, testFile):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, gbdbPackage):
		default:
			return filepath.Base(fmt.Sprintf("%s:%d", file, line))
		}
	}
	return ""
}

func (l *Logger) setLogger() *gblog.Logger {
	var (
		loggerInstanceKey = fmt.Sprintf(`Logger Of Database:%s`, l.group)
	)
	return instances.GetOrSetFuncLock(loggerInstanceKey, func() interface{} {
		logger := l.logger.Clone()
		if l.LogCat != "" {
			logger = logger.Cat(l.LogCat)
		}
		logger.SetStdoutPrint(l.LogStdout)
		return logger
	}).(*gblog.Logger)
}
