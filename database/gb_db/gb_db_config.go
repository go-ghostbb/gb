package gbdb

import (
	"context"
	"ghostbb.io/gb/internal/intlog"
	gblog "ghostbb.io/gb/os/gb_log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type DatabaseConfig struct {
	instance string

	Link          string `json:"lint"`
	Host          string `json:"host"`          // Host of server, ip or domain like: 127.0.0.1, localhost
	Instance      string `json:"instance"`      // instance, if connecting to an instance instead of a port
	Port          string `json:"port"`          // Port, it's commonly 3306.
	User          string `json:"user"`          // Authentication username.
	Pass          string `json:"pass"`          // Authentication password.
	Name          string `json:"name"`          // Default used database name.
	Type          string `json:"type"`          // Database type: mysql, sqlite, mssql, pgsql, oracle.
	Extra         string `json:"extra"`         // (Optional) Extra configuration according the registered third-party database driver.
	Charset       string `json:"charset"`       // (Optional, "utf8" in default) Custom charset when operating on database.
	Protocol      string `json:"protocol"`      // (Optional, "tcp" in default) See net.Dial for more information which networks are available.
	Timezone      string `json:"timezone"`      // (Optional) Sets the time zone for displaying and interpreting time stamps.
	TablePrefix   string `json:"tablePrefix"`   // (Optional) Table name prefix.
	DryRun        bool   `json:"dryRun"`        // (Optional) Dry run, which does SELECT but no INSERT/UPDATE/DELETE statements.
	SingularTable bool   `json:"singularTable"` // (Optional) Use singular table name, table for `User` would be `user` with this option enabled
	MaxIdle       int    `json:"maxIdle"`       // (Optional) Max idle connection configuration for underlying connection pool.
	MaxOpen       int    `json:"maxOpen"`       // (Optional) Max open connection configuration for underlying connection pool.

	// ======================================================================================================
	// Logging.
	// ======================================================================================================
	Terminal          bool          `json:"terminal"`          // Terminal stdout, default true.
	Logger            *gblog.Logger `json:"logger"`            // Logger specifies the logger for server.
	SlowThreshold     time.Duration `json:"slowThreshold"`     // Slow threshold, default 200ms.
	RecordNotFoundErr bool          `json:"recordNotFoundErr"` // Ignore record not found error.
	LogPath           string        `json:"logPath"`           // LogPath specifies the directory for storing logging files.
	LogLevel          string        `json:"logLevel"`          // LogLevel specifies the logging level for logger.
	LogStdout         bool          `json:"logStdout"`         // LogStdout specifies whether printing logging content to stdout.
	ErrorStack        bool          `json:"errorStack"`        // ErrorStack specifies whether logging stack information when error.
	ErrorLogEnabled   bool          `json:"errorLogEnabled"`   // ErrorLogEnabled enables error logging content to files.
	ErrorLogPattern   string        `json:"errorLogPattern"`   // ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	WarnLogEnabled    bool          `json:"warnLogEnabled"`    // WarnLogEnabled enables warn logging content to files.
	WarnLogPattern    string        `json:"warnLogPattern"`    // WarnLogPattern specifies the warn log file pattern like: warn-{Ymd}.log
	AccessLogEnabled  bool          `json:"accessLogEnabled"`  // AccessLogEnabled enables access logging content to files.
	AccessLogPattern  string        `json:"accessLogPattern"`  // AccessLogPattern specifies the access log file pattern like: access-{Ymd}.log

	// ======================================================================================================
	// Cluster.
	// ======================================================================================================
	Role string `json:"role"` // (Optional, "master" in default) Node role, used for master-slave mode: master, slave.
}

func NewConfig() DatabaseConfig {
	return DatabaseConfig{
		Charset:           defaultCharset,
		Protocol:          defaultProtocol,
		MaxIdle:           10,
		Terminal:          true,
		Logger:            gblog.New(),
		SlowThreshold:     200 * time.Millisecond,
		RecordNotFoundErr: true,
		LogLevel:          "all",
		LogStdout:         false,
		ErrorStack:        true,
		ErrorLogEnabled:   true,
		ErrorLogPattern:   "error-{Ymd}.log",
		WarnLogEnabled:    true,
		WarnLogPattern:    "warn-{Ymd}.log",
		AccessLogEnabled:  false,
		AccessLogPattern:  "access-{Ymd}.log",
	}
}

func (d *DB) setConfig(config DatabaseConfig) error {
	d.config = config

	// Logging.
	if d.config.LogPath != "" && d.config.LogPath != d.config.Logger.GetPath() {
		if err := d.config.Logger.SetPath(d.config.LogPath); err != nil {
			return err
		}
	}
	if err := d.config.Logger.SetLevelStr(d.config.LogLevel); err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	intlog.Printf(context.TODO(), "SetConfig: %+v", d.config)
	return nil
}

func (d *DB) Logger() *gblog.Logger {
	return d.config.Logger
}

func (d *DatabaseConfig) GormConfig() *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   d.TablePrefix,
			SingularTable: d.SingularTable,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		DryRun:                                   d.DryRun,
	}
	config.Logger = d.newDBLogger()
	return config
}

func (d *DatabaseConfig) IsErrorLogEnabled() bool {
	return d.ErrorLogEnabled
}

func (d *DatabaseConfig) IsAccessLogEnabled() bool {
	return d.AccessLogEnabled
}

func (d *DatabaseConfig) IsWarnLogEnabled() bool {
	return d.WarnLogEnabled
}
