package dbex

import (
	"time"
	"os"
	"fmt"
	goLog "log"
	"runtime"
)

type Level uint8

const (
	LevelDebug   Level = 0
	LevelInfo    Level = 1
	LevelWarning Level = 2
	LevelError   Level = 3
	LevelPanic   Level = 4
)

type Logger struct {
	level         Level
	logFileName   string
	activeLog     bool
	activeLogFile bool
}

func newLogger(fileNamePrefix string) (*Logger, error) {
	l := &Logger{
		level:         LevelWarning,
		logFileName:   fileNamePrefix,
		activeLog:     true,
		activeLogFile: true,
	}
	return l, nil
}
func (l *Logger) SetLogFileName(name string) {
	l.logFileName = name
}
func (l *Logger) SetLevel(level Level) {
	l.level = level
}
func (l *Logger) ActiveLog(active bool) {
	l.activeLog = active
}
func (l *Logger) ActiveLogFile(active bool) {
	l.activeLogFile = active
}
func (l *Logger) Log(logLevel Level, msg string) {

	if !l.activeLog {
		return
	}

	prefix := "[Info]"
	switch logLevel {
	case LevelDebug:
		prefix = "[Debug]"
	case LevelWarning:
		prefix = "[warning]"
	case LevelError:
		prefix = "[Error]"
	case LevelPanic:
		prefix = "[Panic]"

	}

	_, file, line, _ := runtime.Caller(1)

	goLog.Printf("%s%s:%d %v", prefix, file, line, msg)
}

func (l *Logger) LogFile(logLevel Level, v ...interface{}) {

	if !l.activeLogFile {
		return
	}

	if logLevel >= l.level {

		prefix := "[Info]"
		switch logLevel {
		case LevelDebug:
			prefix = "[Debug]"
		case LevelWarning:
			prefix = "[warning]"
		case LevelError:
			prefix = "[Error]"
		case LevelPanic:
			prefix = "[Panic]"

		}
		filePath := l.logFileName + time.Now().Format("20060102") + ".log"
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err == nil {
			lr := goLog.New(f, prefix, goLog.LstdFlags|goLog.Lshortfile)
			lr.Output(5, fmt.Sprintln(v))

		}
		defer f.Close()

	}
}
