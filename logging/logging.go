//
//
//

// Package logging, custom log
package logging

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"
	"strings"
)

//	Interface: logger
//
// # 機能
//
//	ログレベルに応じたログを出力する
type Logging interface {
	// Debug レベルのログを出力する
	Debug(format string, args ...interface{})
	// Info レベルのログを出力する
	Info(format string, args ...interface{})
	// Warning レベルのログを出力する
	Warn(format string, args ...interface{})
	// Error レベルのログを出力する
	Error(format string, args ...interface{})
	// Fatal レベルのログを出力する
	Fatal(format string, args ...interface{})
}

// Struct: Logger
//
// # ロガーの構造体
type logger struct {
	stderr *stdlog.Logger
	stdout *stdlog.Logger
	level  int
}

// ログレベル
type Level int

// ログレベルの定数
const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG
)

// Function: NewLogger()
//
// # 新規ロガーを取得する
func NewLogger() *logger {
	return &logger{
		stdout: stdlog.New(os.Stdout, "", 0),
		stderr: stdlog.New(os.Stderr, "", 0),
		level:  SetLogLevel(),
	}
}

// Function: SetlogLevel()
//
// # システム環境変数よりログレベルを設定する
func SetLogLevel() int {
	level := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch level {
	case "DEBUG":
		return int(DEBUG)
	case "INFO":
		return int(INFO)
	case "WARN":
		return int(WARN)
	case "ERROR":
		return int(ERROR)
	case "FATAL":
		return int(FATAL)
	default:
		return int(INFO)
	}
}

// Method: Level.String()
//
// ログレベル定数のStriggerメソッド
func (l Level) String() string {
	switch l {
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "INFO"
	}
}

// Method: logger.Fatail()
//
// Fatail レベルのログを出力する
func (l *logger) Fatal(format string, args ...interface{}) {
	if l.level >= int(FATAL) {
		l.stderr.SetPrefix(fmt.Sprintf("[!!!%s!!!] ", FATAL))
		l.stderr.SetFlags(stdlog.Ldate | stdlog.Ltime)

		_, file, line, ok := runtime.Caller(1)
		if ok {
			caller := fmt.Sprintf("@%s:%d: ", file, line)
			l.stderr.Panicf(caller+format, args...)
		} else {
			l.stderr.Panicf(format, args...)
		}
	}
}

// Method: logger.Error()
//
// Error レベルのログを出力する
func (l *logger) Error(format string, args ...interface{}) {
	if l.level >= int(ERROR) {
		l.stderr.SetPrefix(fmt.Sprintf("[%s] ", ERROR))
		l.stderr.SetFlags(stdlog.Ldate | stdlog.Ltime)

		_, file, line, ok := runtime.Caller(1)
		if ok {
			caller := fmt.Sprintf("@%s:%d: ", file, line)
			l.stderr.Printf(caller+format, args...)
		} else {
			l.stderr.Printf(format, args...)
		}
	}
}

// Method: logger.Warn()
//
// Warn レベルのログを出力する
func (l *logger) Warn(format string, args ...interface{}) {
	if l.level >= int(WARN) {
		l.stderr.SetPrefix(fmt.Sprintf("[%s] ", WARN))
		l.stderr.SetFlags(stdlog.Ldate | stdlog.Ltime)
		l.stderr.Printf(format, args...)
	}
}

// Method: logger.Info()
//
// Info レベルのログを出力する
func (l *logger) Info(format string, args ...interface{}) {
	if l.level >= int(INFO) {
		l.stdout.SetPrefix(fmt.Sprintf("[%s] ", INFO))
		l.stdout.SetFlags(stdlog.Ldate | stdlog.Ltime)
		l.stdout.Printf(format, args...)
	}
}

// Method: logger.Debug()
//
// Debug レベルのログを出力する
func (l *logger) Debug(format string, args ...interface{}) {
	if l.level >= int(DEBUG) {
		l.stdout.SetPrefix(fmt.Sprintf("[%s] ", DEBUG))
		l.stdout.SetFlags(stdlog.Ldate | stdlog.Ltime)

		_, file, line, ok := runtime.Caller(1)
		if ok {
			caller := fmt.Sprintf("@%s:%d: ", file, line)
			l.stdout.Printf(caller+format, args...)
		} else {
			l.stdout.Printf(format, args...)
		}
	}
}
