package gblog

import (
	"context"
	"fmt"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbcompress "ghostbb.io/gb/encoding/gb_compress"
	"ghostbb.io/gb/internal/intlog"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbmlock "ghostbb.io/gb/os/gb_mlock"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbregex "ghostbb.io/gb/text/gb_regex"
	"strings"
	"time"
)

const (
	memoryLockPrefixForRotating = "gblog.rotateChecksTimely:"
)

// rotateFileBySize rotates the current logging file according to the
// configured rotation size.
func (l *Logger) rotateFileBySize(ctx context.Context, now time.Time) {
	if l.config.RotateSize <= 0 {
		return
	}
	if err := l.doRotateFile(ctx, l.getFilePath(now)); err != nil {
		// panic(err)
		intlog.Errorf(ctx, `%+v`, err)
	}
}

// doRotateFile rotates the given logging file.
func (l *Logger) doRotateFile(ctx context.Context, filePath string) error {
	memoryLockKey := "gblog.doRotateFile:" + filePath
	if !gbmlock.TryLock(memoryLockKey) {
		return nil
	}
	defer gbmlock.Unlock(memoryLockKey)

	intlog.PrintFunc(ctx, func() string {
		return fmt.Sprintf(`start rotating file by size: %s, file: %s`, gbfile.SizeFormat(filePath), filePath)
	})
	defer intlog.PrintFunc(ctx, func() string {
		return fmt.Sprintf(`done rotating file by size: %s, size: %s`, gbfile.SizeFormat(filePath), filePath)
	})

	// No backups, it then just removes the current logging file.
	if l.config.RotateBackupLimit == 0 {
		if err := gbfile.Remove(filePath); err != nil {
			return err
		}
		intlog.Printf(
			ctx,
			`%d size exceeds, no backups set, remove original logging file: %s`,
			l.config.RotateSize, filePath,
		)
		return nil
	}
	// Else it creates new backup files.
	var (
		dirPath     = gbfile.Dir(filePath)
		fileName    = gbfile.Name(filePath)
		fileExtName = gbfile.ExtName(filePath)
		newFilePath = ""
	)
	// Rename the logging file by adding extra datetime information to microseconds, like:
	// access.log          -> access.20200326101301899002.log
	// access.20200326.log -> access.20200326.20200326101301899002.log
	for {
		var (
			now   = gbtime.Now()
			micro = now.Microsecond() % 1000
		)
		if micro == 0 {
			micro = 101
		} else {
			for micro < 100 {
				micro *= 10
			}
		}
		newFilePath = gbfile.Join(
			dirPath,
			fmt.Sprintf(
				`%s.%s%d.%s`,
				fileName, now.Format("YmdHisu"), micro, fileExtName,
			),
		)
		if !gbfile.Exists(newFilePath) {
			break
		} else {
			intlog.Printf(ctx, `rotation file exists, continue: %s`, newFilePath)
		}
	}
	intlog.Printf(ctx, "rotating file by size from %s to %s", filePath, newFilePath)
	if err := gbfile.Rename(filePath, newFilePath); err != nil {
		return err
	}
	return nil
}

// rotateChecksTimely timely checks the backups expiration and the compression.
func (l *Logger) rotateChecksTimely(ctx context.Context) {
	defer gbtimer.AddOnce(ctx, l.config.RotateCheckInterval, l.rotateChecksTimely)

	// Checks whether file rotation not enabled.
	if l.config.RotateSize <= 0 && l.config.RotateExpire == 0 {
		intlog.Printf(
			ctx,
			"logging rotation ignore checks: RotateSize: %d, RotateExpire: %s",
			l.config.RotateSize, l.config.RotateExpire.String(),
		)
		return
	}

	// It here uses memory lock to guarantee the concurrent safety.
	memoryLockKey := memoryLockPrefixForRotating + l.config.Path
	if !gbmlock.TryLock(memoryLockKey) {
		return
	}
	defer gbmlock.Unlock(memoryLockKey)

	var (
		now        = time.Now()
		pattern    = "*.log, *.gz"
		files, err = gbfile.ScanDirFile(l.config.Path, pattern, true)
	)
	if err != nil {
		intlog.Errorf(ctx, `%+v`, err)
	}
	intlog.Printf(ctx, "logging rotation start checks: %+v", files)
	// get file name regex pattern
	// access-{y-m-d}-test.log => access-$-test.log => access-\$-test\.log => access-(.+?)-test\.log
	fileNameRegexPattern, _ := gbregex.ReplaceString(`{.+?}`, "$", l.config.File)
	fileNameRegexPattern = gbregex.Quote(fileNameRegexPattern)
	fileNameRegexPattern = strings.ReplaceAll(fileNameRegexPattern, "\\$", "(.+?)")
	// =============================================================
	// Rotation of expired file checks.
	// =============================================================
	if l.config.RotateExpire > 0 {
		var (
			mtime         time.Time
			subDuration   time.Duration
			expireRotated bool
		)
		for _, file := range files {
			// ignore backup file
			if gbregex.IsMatchString(`.+\.\d{20}\.log`, gbfile.Basename(file)) || gbfile.ExtName(file) == "gz" {
				continue
			}
			// ignore not matching file
			if !gbregex.IsMatchString(fileNameRegexPattern, file) {
				continue
			}
			mtime = gbfile.MTime(file)
			subDuration = now.Sub(mtime)
			if subDuration > l.config.RotateExpire {
				func() {
					memoryLockFileKey := memoryLockPrefixForPrintingToFile + file
					if !gbmlock.TryLock(memoryLockFileKey) {
						return
					}
					defer gbmlock.Unlock(memoryLockFileKey)
					expireRotated = true
					intlog.Printf(
						ctx,
						`%v - %v = %v > %v, rotation expire logging file: %s`,
						now, mtime, subDuration, l.config.RotateExpire, file,
					)
					if err = l.doRotateFile(ctx, file); err != nil {
						intlog.Errorf(ctx, `%+v`, err)
					}
				}()
			}
		}
		if expireRotated {
			// Update the files array.
			files, err = gbfile.ScanDirFile(l.config.Path, pattern, true)
			if err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		}
	}

	// =============================================================
	// Rotated file compression.
	// =============================================================
	needCompressFileArray := gbarray.NewStrArray()
	if l.config.RotateBackupCompress > 0 {
		for _, file := range files {
			// Eg: access.20200326101301899002.log.gz
			if gbfile.ExtName(file) == "gz" {
				continue
			}
			// ignore not matching file
			originalLoggingFilePath, _ := gbregex.ReplaceString(`\.\d{20}`, "", file)
			if !gbregex.IsMatchString(fileNameRegexPattern, originalLoggingFilePath) {
				continue
			}
			// Eg:
			// access.20200326101301899002.log
			if gbregex.IsMatchString(`.+\.\d{20}\.log`, gbfile.Basename(file)) {
				needCompressFileArray.Append(file)
			}
		}
		if needCompressFileArray.Len() > 0 {
			needCompressFileArray.Iterator(func(_ int, path string) bool {
				err := gbcompress.GzipFile(path, path+".gz")
				if err == nil {
					intlog.Printf(ctx, `compressed done, remove original logging file: %s`, path)
					if err = gbfile.Remove(path); err != nil {
						intlog.Print(ctx, err)
					}
				} else {
					intlog.Print(ctx, err)
				}
				return true
			})
			// Update the files array.
			files, err = gbfile.ScanDirFile(l.config.Path, pattern, true)
			if err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		}
	}

	// =============================================================
	// Backups count limitation and expiration checks.
	// =============================================================
	backupFiles := gbarray.NewSortedArray(func(a, b interface{}) int {
		// Sorted by rotated/backup file mtime.
		// The older rotated/backup file is put in the head of array.
		var (
			file1  = a.(string)
			file2  = b.(string)
			result = gbfile.MTimestampMilli(file1) - gbfile.MTimestampMilli(file2)
		)
		if result <= 0 {
			return -1
		}
		return 1
	})
	if l.config.RotateBackupLimit > 0 || l.config.RotateBackupExpire > 0 {
		for _, file := range files {
			// ignore not matching file
			originalLoggingFilePath, _ := gbregex.ReplaceString(`\.\d{20}`, "", file)
			if !gbregex.IsMatchString(fileNameRegexPattern, originalLoggingFilePath) {
				continue
			}
			if gbregex.IsMatchString(`.+\.\d{20}\.log`, gbfile.Basename(file)) {
				backupFiles.Add(file)
			}
		}
		intlog.Printf(ctx, `calculated backup files array: %+v`, backupFiles)
		diff := backupFiles.Len() - l.config.RotateBackupLimit
		for i := 0; i < diff; i++ {
			path, _ := backupFiles.PopLeft()
			intlog.Printf(ctx, `remove exceeded backup limit file: %s`, path)
			if err := gbfile.Remove(path.(string)); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		}
		// Backups expiration checking.
		if l.config.RotateBackupExpire > 0 {
			var (
				mtime       time.Time
				subDuration time.Duration
			)
			backupFiles.Iterator(func(_ int, v interface{}) bool {
				path := v.(string)
				mtime = gbfile.MTime(path)
				subDuration = now.Sub(mtime)
				if subDuration > l.config.RotateBackupExpire {
					intlog.Printf(
						ctx,
						`%v - %v = %v > %v, remove expired backup file: %s`,
						now, mtime, subDuration, l.config.RotateBackupExpire, path,
					)
					if err := gbfile.Remove(path); err != nil {
						intlog.Errorf(ctx, `%+v`, err)
					}
					return true
				} else {
					return false
				}
			})
		}
	}
}
