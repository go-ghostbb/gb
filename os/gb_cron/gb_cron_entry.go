package gbcron

import (
	"context"
	"fmt"
	gbtype "ghostbb.io/gb/container/gb_type"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gblog "ghostbb.io/gb/os/gb_log"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"reflect"
	"runtime"
	"time"
)

// JobFunc is the timing called job function in cron.
type JobFunc = gbtimer.JobFunc

// Entry is timing task entry.
type Entry struct {
	cron       *Cron          // Cron object belonged to.
	timerEntry *gbtimer.Entry // Associated timer Entry.
	schedule   *cronSchedule  // Timed schedule object.
	jobName    string         // Callback function name(address info).
	times      *gbtype.Int    // Running times limit.
	infinite   *gbtype.Bool   // No times limit.
	Name       string         // Entry name.
	Job        JobFunc        `json:"-"` // Callback function.
	Time       time.Time      // Registered time.
}

type doAddEntryInput struct {
	Name        string          // Name names this entry for manual control.
	Job         JobFunc         // Job is the callback function for timed task execution.
	Ctx         context.Context // The context for the job.
	Times       int             // Times specifies the running limit times for the entry.
	Pattern     string          // Pattern is the crontab style string for scheduler.
	IsSingleton bool            // Singleton specifies whether timed task executing in singleton mode.
	Infinite    bool            // Infinite specifies whether this entry is running with no times limit.
}

// doAddEntry creates and returns a new Entry object.
func (c *Cron) doAddEntry(in doAddEntryInput) (*Entry, error) {
	if in.Name != "" {
		if c.Search(in.Name) != nil {
			return nil, gberror.NewCodef(gbcode.CodeInvalidOperation, `cron job "%s" already exists`, in.Name)
		}
	}
	schedule, err := newSchedule(in.Pattern)
	if err != nil {
		return nil, err
	}
	// No limit for `times`, for timer checking scheduling every second.
	entry := &Entry{
		cron:     c,
		schedule: schedule,
		jobName:  runtime.FuncForPC(reflect.ValueOf(in.Job).Pointer()).Name(),
		times:    gbtype.NewInt(in.Times),
		infinite: gbtype.NewBool(in.Infinite),
		Job:      in.Job,
		Time:     time.Now(),
	}
	if in.Name != "" {
		entry.Name = in.Name
	} else {
		entry.Name = "cron-" + gbconv.String(c.idGen.Add(1))
	}
	// When you add a scheduled task, you cannot allow it to run.
	// It cannot start running when added to timer.
	// It should start running after the entry is added to the Cron entries map, to avoid the task
	// from running during adding where the entries do not have the entry information, which might cause panic.
	entry.timerEntry = gbtimer.AddEntry(
		in.Ctx,
		time.Second,
		entry.checkAndRun,
		in.IsSingleton,
		-1,
		gbtimer.StatusStopped,
	)
	c.entries.Set(entry.Name, entry)
	entry.timerEntry.Start()
	return entry, nil
}

// IsSingleton return whether this entry is a singleton timed task.
func (entry *Entry) IsSingleton() bool {
	return entry.timerEntry.IsSingleton()
}

// SetSingleton sets the entry running in singleton mode.
func (entry *Entry) SetSingleton(enabled bool) {
	entry.timerEntry.SetSingleton(enabled)
}

// SetTimes sets the times which the entry can run.
func (entry *Entry) SetTimes(times int) {
	entry.times.Set(times)
	entry.infinite.Set(false)
}

// Status returns the status of entry.
func (entry *Entry) Status() int {
	return entry.timerEntry.Status()
}

// SetStatus sets the status of the entry.
func (entry *Entry) SetStatus(status int) int {
	return entry.timerEntry.SetStatus(status)
}

// Start starts running the entry.
func (entry *Entry) Start() {
	entry.timerEntry.Start()
}

// Stop stops running the entry.
func (entry *Entry) Stop() {
	entry.timerEntry.Stop()
}

// Close stops and removes the entry from cron.
func (entry *Entry) Close() {
	entry.cron.entries.Remove(entry.Name)
	entry.timerEntry.Close()
}

// checkAndRun is the core timing task check logic.
func (entry *Entry) checkAndRun(ctx context.Context) {
	currentTime := time.Now()
	if !entry.schedule.checkMeetAndUpdateLastSeconds(ctx, currentTime) {
		return
	}
	switch entry.cron.status.Val() {
	case StatusStopped:
		return

	case StatusClosed:
		entry.logDebugf(ctx, `cron job "%s" is removed`, entry.getJobNameWithPattern())
		entry.Close()

	case StatusReady, StatusRunning:
		defer func() {
			if exception := recover(); exception != nil {
				// Exception caught, it logs the error content to logger in default behavior.
				entry.logErrorf(ctx,
					`cron job "%s(%s)" end with error: %+v`,
					entry.jobName, entry.schedule.pattern, exception,
				)
			} else {
				entry.logDebugf(ctx, `cron job "%s" ends`, entry.getJobNameWithPattern())
			}
			if entry.timerEntry.Status() == StatusClosed {
				entry.Close()
			}
		}()

		// Running times check.
		if !entry.infinite.Val() {
			times := entry.times.Add(-1)
			if times <= 0 {
				if entry.timerEntry.SetStatus(StatusClosed) == StatusClosed || times < 0 {
					return
				}
			}
		}
		entry.logDebugf(ctx, `cron job "%s" starts`, entry.getJobNameWithPattern())
		entry.Job(ctx)
	}
}

func (entry *Entry) getJobNameWithPattern() string {
	return fmt.Sprintf(`%s(%s)`, entry.jobName, entry.schedule.pattern)
}

func (entry *Entry) logDebugf(ctx context.Context, format string, v ...interface{}) {
	if logger := entry.cron.GetLogger(); logger != nil {
		logger.Debugf(ctx, format, v...)
	}
}

func (entry *Entry) logErrorf(ctx context.Context, format string, v ...interface{}) {
	logger := entry.cron.GetLogger()
	if logger == nil {
		logger = gblog.DefaultLogger()
	}
	logger.Errorf(ctx, format, v...)
}
