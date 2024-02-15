package gbdb

import (
	"context"
	"fmt"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/instance"
	gblog "ghostbb.io/gb/os/gb_log"
	gbstr "ghostbb.io/gb/text/gb_str"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
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

func (d *DatabaseConfig) newDBLogger() *dbLogger {
	return &dbLogger{
		DatabaseConfig: d,
	}
}

type dbLogger struct {
	*DatabaseConfig
}

func (d *dbLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return d
}

func (d *dbLogger) Info(ctx context.Context, str string, args ...interface{}) {
	logger := d.setLogger()
	logger.Infof(ctx, d.instance+" "+str, args...)
}

func (d *dbLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	logger := d.setLogger()
	logger.Warningf(ctx, d.instance+" "+str, args...)
}

func (d *dbLogger) Error(ctx context.Context, str string, args ...interface{}) {
	logger := d.setLogger()
	logger.Errorf(ctx, d.instance+" "+str, args...)
}

func (d *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	var (
		elapsed   = time.Since(begin)
		logger    = d.setLogger()
		track     = d.pathTrack(ctx)
		sql, rows = fc()
		msg       = fmt.Sprintf(`%s 【%s】｜%10v｜%10s`,
			track,
			d.instance,
			elapsed,
			fmt.Sprintf("rows:%d", rows),
		)
		elapsedColor = d.ElapsedColor(elapsed)
		resetColor   = reset
	)
	sql = gbstr.Replace(sql, `/`, "", -1)
	msg += "\n    " + sql + "\n"

	switch {
	case err != nil && (d.RecordNotFoundErr || !gberror.Is(err, gorm.ErrRecordNotFound)):
		logger.Error(ctx, msg, map[string]interface{}{"error": err.Error()})
	case d.SlowThreshold != 0 && elapsed > d.SlowThreshold:
		logger.Warning(ctx, msg)
	default:
		logger.Info(ctx, msg)
	}

	if d.Terminal {
		fmt.Printf("%s [ORM] %s ｜%s %15v %s｜%s %-8s %s｜%-10s｜ %s \n",
			time.Now().Format("2006/01/02 15:04:05"),
			track,
			elapsedColor, elapsed, resetColor,
			magenta, fmt.Sprintf("rows:%d", rows), reset,
			d.instance,
			sql,
		)
	}
}

func (d *dbLogger) ElapsedColor(elapsed time.Duration) string {
	if elapsed > d.SlowThreshold {
		return red
	}
	return green
}

const (
	gormPackage = "gorm.io"
	gbdbPackage = "database/gb_db"
	testFile    = "_test.go"
	genFile     = ".gen.go"
)

func (d *dbLogger) pathTrack(ctx context.Context) string {
	// Caller lookup
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		switch {
		case !ok:
		case gbstr.HasSuffix(file, testFile):
		case gbstr.HasSuffix(file, genFile):
		case gbstr.Contains(file, gormPackage):
		case gbstr.Contains(file, gbdbPackage):
		default:
			return filepath.Base(fmt.Sprintf("%s:%d", file, line))
		}
	}
	return ""
}

func (d *dbLogger) setLogger() *gblog.Logger {
	var (
		loggerInstanceKey = fmt.Sprintf(`Logger Of Database:%s`, d.instance)
	)
	return instance.GetOrSetFuncLock(loggerInstanceKey, func() interface{} {
		logger := d.Logger.Clone()
		logger.SetFile(d.AccessLogPattern)
		logger.SetStdoutPrint(d.LogStdout)
		logger.SetLevelPrint(false)
		return logger
	}).(*gblog.Logger)
}
