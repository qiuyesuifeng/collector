// refer to github.com/ngaut/logging
// high level log wrapper, so it can output different log based on level
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

const (
	Ldate         = log.Ldate
	Llongfile     = log.Llongfile
	Lmicroseconds = log.Lmicroseconds
	Lshortfile    = log.Lshortfile
	LstdFlags     = log.LstdFlags
	Ltime         = log.Ltime
)

type (
	LogLevel int
	LogType  int
)

const (
	LOG_FATAL   = LogType(0x1)
	LOG_ERROR   = LogType(0x2)
	LOG_WARNING = LogType(0x4)
	LOG_INFO    = LogType(0x8)
	LOG_DEBUG   = LogType(0x10)
)

const (
	LOG_LEVEL_NONE  = LogLevel(0x0)
	LOG_LEVEL_FATAL = LOG_LEVEL_NONE | LogLevel(LOG_FATAL)
	LOG_LEVEL_ERROR = LOG_LEVEL_FATAL | LogLevel(LOG_ERROR)
	LOG_LEVEL_WARN  = LOG_LEVEL_ERROR | LogLevel(LOG_WARNING)
	LOG_LEVEL_INFO  = LOG_LEVEL_WARN | LogLevel(LOG_INFO)
	LOG_LEVEL_DEBUG = LOG_LEVEL_INFO | LogLevel(LOG_DEBUG)
	LOG_LEVEL_ALL   = LOG_LEVEL_DEBUG
)

const FORMAT_TIME_DAY string = "20060102"
const FORMAT_TIME_HOUR string = "2006010215"

var Log *Logger = New()

func init() {
	SetFlags(Ldate | Ltime | Lshortfile)
}

func CrashLog(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
	} else {
		syscall.Dup2(int(f.Fd()), 2)
	}
}

func SetLevel(level LogLevel) {
	Log.SetLevel(level)
}
func GetLogLevel() LogLevel {
	return Log.level
}

func SetOutput(out io.Writer) {
	Log.SetOutput(out)
}

func SetOutputByName(path string) error {
	return Log.SetOutputByName(path)
}

func SetFlags(flags int) {
	Log.Log.SetFlags(flags)
}

func Info(format string, v ...interface{}) {
	Log.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	Log.Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	Log.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	Log.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	Log.Fatal(format, v...)
}

func SetLevelByString(level string) {
	Log.SetLevelByString(level)
}

func SetHighlighting(highlighting bool) {
	Log.SetHighlighting(highlighting)
}

func SetRotateByDay() {
	Log.SetRotateByDay()
}

func SetRotateByHour() {
	Log.SetRotateByHour()
}

type Logger struct {
	Log          *log.Logger
	level        LogLevel
	highlighting bool

	dailyRolling bool
	hourRolling  bool

	fileName  string
	logSuffix string
	fd        *os.File

	lock sync.Mutex
}

func (l *Logger) SetHighlighting(highlighting bool) {
	l.highlighting = highlighting
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetLevelByString(level string) {
	l.level = StringToLogLevel(level)
}

func (l *Logger) SetRotateByDay() {
	l.dailyRolling = true
	l.logSuffix = genDayTime(time.Now())
}

func (l *Logger) SetRotateByHour() {
	l.hourRolling = true
	l.logSuffix = genHourTime(time.Now())
}

func (l *Logger) rotate() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	var suffix string
	if l.dailyRolling {
		suffix = genDayTime(time.Now())
	} else if l.hourRolling {
		suffix = genHourTime(time.Now())
	} else {
		return nil
	}

	// Notice: if suffix is not equal to l.LogSuffix, then rotate
	if suffix != l.logSuffix {
		err := l.doRotate(suffix)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Logger) doRotate(suffix string) error {
	lastFileName := l.fileName + "." + l.logSuffix
	err := os.Rename(l.fileName, lastFileName)
	if err != nil {
		return err
	}

	// Notice: Not check error, is this ok?
	l.fd.Close()

	err = l.SetOutputByName(l.fileName)
	if err != nil {
		return err
	}

	l.logSuffix = suffix

	return nil
}

func (l *Logger) SetOutput(out io.Writer) {
	l.Log = log.New(out, l.Log.Prefix(), l.Log.Flags())
}

func (l *Logger) SetOutputByName(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l.SetOutput(f)

	l.fileName = path
	l.fd = f

	return err
}

func (l *Logger) log(t LogType, v ...interface{}) {
	if l.level|LogLevel(t) != l.level {
		return
	}

	err := l.rotate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	v1 := make([]interface{}, len(v)+2)
	logStr, logColor := LogTypeToString(t)
	if l.highlighting {
		v1[0] = "\033" + logColor + "m[" + logStr + "]"
		copy(v1[1:], v)
		v1[len(v)+1] = "\033[0m"
	} else {
		v1[0] = "[" + logStr + "]"
		copy(v1[1:], v)
		v1[len(v)+1] = ""
	}

	s := fmt.Sprintln(v1...)
	l.Log.Output(4, s)
}

func (l *Logger) logf(t LogType, format string, v ...interface{}) {
	if l.level|LogLevel(t) != l.level {
		return
	}

	err := l.rotate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	logStr, logColor := LogTypeToString(t)
	var s string
	if l.highlighting {
		s = "\033" + logColor + "m[" + logStr + "] " + fmt.Sprintf(format, v...) + "\033[0m"
	} else {
		s = "[" + logStr + "] " + fmt.Sprintf(format, v...)
	}
	l.Log.Output(4, s)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logf(LOG_FATAL, format, v...)
	os.Exit(-1)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logf(LOG_ERROR, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logf(LOG_WARNING, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.logf(LOG_DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logf(LOG_INFO, format, v...)
}

// create a output function
func (l *Logger) Output(maxdepth int, s string) error {
	// 0 --> debug
	// 1 --> info
	// 2 --> info
	// 3 --> error
	if maxdepth == 0 {
		l.Debug(s)
	} else if maxdepth == 1 {
		l.Info(s)
	} else if maxdepth == 2 {
		l.Info(s)
	} else if maxdepth == 3 {
		l.Error(s)
	}

	return nil
}

func StringToLogLevel(level string) LogLevel {
	switch level {
	case "fatal":
		return LOG_LEVEL_FATAL
	case "error":
		return LOG_LEVEL_ERROR
	case "warn":
		return LOG_LEVEL_WARN
	case "warning":
		return LOG_LEVEL_WARN
	case "debug":
		return LOG_LEVEL_DEBUG
	case "info":
		return LOG_LEVEL_INFO
	}
	return LOG_LEVEL_ALL
}

func LogTypeToString(t LogType) (string, string) {
	switch t {
	case LOG_FATAL:
		return "fatal", "[0;31"
	case LOG_ERROR:
		return "error", "[0;31"
	case LOG_WARNING:
		return "warning", "[0;33"
	case LOG_DEBUG:
		return "debug", "[0;36"
	case LOG_INFO:
		return "info", "[0;37"
	}
	return "unknown", "[0;37"
}

func genDayTime(t time.Time) string {
	return t.Format(FORMAT_TIME_DAY)
}

func genHourTime(t time.Time) string {
	return t.Format(FORMAT_TIME_HOUR)
}

func New() *Logger {
	return Newlogger(os.Stdout, "")
}

func Newlogger(w io.Writer, prefix string) *Logger {
	return &Logger{Log: log.New(w, prefix, LstdFlags), level: LOG_LEVEL_ALL, highlighting: true}
}
