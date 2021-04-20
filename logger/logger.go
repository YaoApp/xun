// Fork from https://github.com/gin-gonic/gin/blob/master/logger.go

package logger

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mattn/go-isatty"
	"github.com/yaoapp/xun/global"
)

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	consoleColorMode = disableColor
}

// ForceConsoleColor force color output in the console.
func ForceConsoleColor() {
	consoleColorMode = forceColor
}

func init() {
	DefaultLogger = New(
		global.DefaultWriter,
		DefaultFormatter,
		DefaultLevel,
	)
	DefaultErrorLogger = New(
		global.DefaultErrorWriter,
		DefaultFormatter,
		DefaultLevel,
	)
}

// New create a new logger instance
func New(out io.Writer, f Formatter, level LogLevel, notlogged ...string) *Logger {
	logger := &Logger{}
	logger.SetOutput(out)
	logger.SetFormatter(f)
	logger.SetNotlogged(notlogged...)
	logger.SetLevel(level)
	return logger
}

// SetDefaultLogger set the default logger instance
func SetDefaultLogger(logger *Logger) {
	DefaultLogger = logger
}

// SetDefaultErrorLogger set the default error logger instance
func SetDefaultErrorLogger(logger *Logger) {
	DefaultErrorLogger = logger
}

// SetDefaultLevel set the default log level
func SetDefaultLevel(level LogLevel) {
	DefaultErrorLogger.SetLevel(level)
	DefaultLogger.SetLevel(level)
}

// SetOutput  set the output writer to the given value
func (logger *Logger) SetOutput(out io.Writer) *Logger {
	isTerm := true
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}
	logger.isTerm = isTerm
	logger.Output = out
	return logger
}

// SetFormatter  set the formatter to the given value
func (logger *Logger) SetFormatter(f Formatter) *Logger {
	logger.Formatter = f
	return logger
}

// SetNotlogged  set the not logged methods
func (logger *Logger) SetNotlogged(notlogged ...string) *Logger {
	logger.SkipMethods = notlogged
	return logger
}

// SetLevel  set the level to the given value
func (logger *Logger) SetLevel(level LogLevel) *Logger {
	logger.Level = level
	return logger
}

// Fatal represents truly catastrophic situations, as far as your application is concerned.
func Fatal(code int, message string) Writer {
	return DefaultErrorLogger.Fatal(code, message)
}

// Fatal represents truly catastrophic situations, as far as your application is concerned.
func (logger *Logger) Fatal(code int, message string) Writer {
	if logger.Level < LevelFatal {
		return Empty{}
	}

	pc, file, line, ok := runtime.Caller(2)
	method := "unknown"
	if ok {
		details := runtime.FuncForPC(pc)
		method = filepath.Base(details.Name())
	}
	return &Content{
		isTerm: logger.isTerm,
		logger: logger,
		Type:   ERROR,
		Method: method,
		Level:  LevelFatal,
		Code:   code,
		File:   file,
		Line:   line,
		Msg:    message,
	}
}

// Error An error is a serious issue and represents the failure of something important going on in your application.
func Error(code int, message string) Writer {
	return DefaultErrorLogger.Error(code, message)
}

// Error An error is a serious issue and represents the failure of something important going on in your application.
func (logger *Logger) Error(code int, message string) Writer {
	if logger.Level < LevelError {
		return Empty{}
	}
	pc, file, line, ok := runtime.Caller(2)
	method := "unknown"
	if ok {
		details := runtime.FuncForPC(pc)
		method = filepath.Base(details.Name())
	}
	return &Content{
		isTerm: logger.isTerm,
		logger: logger,
		Type:   ERROR,
		Method: method,
		Level:  LevelError,
		Code:   code,
		File:   file,
		Line:   line,
		Msg:    message,
	}
}

// Warn Now we're getting into the grayer area of hypotheticals.
func Warn(methodType string, trace ...string) Writer {
	return DefaultLogger.Warn(methodType, trace...)
}

// Warn Now we're getting into the grayer area of hypotheticals.
func (logger *Logger) Warn(methodType string, trace ...string) Writer {
	if logger.Level < LevelWarn {
		return Empty{}
	}
	pc, file, line, ok := runtime.Caller(2)
	method := "unknown"
	if ok {
		details := runtime.FuncForPC(pc)
		method = filepath.Base(details.Name())
	}
	return &Content{
		isTerm: logger.isTerm,
		logger: logger,
		Type:   methodType,
		Method: method,
		Level:  LevelWarn,
		File:   file,
		Line:   line,
		Traces: trace,
	}
}

// Info Finally, we can dial down the stress level.
func Info(methodType string, trace ...string) Writer {
	return DefaultLogger.Info(methodType, trace...)
}

// Info Finally, we can dial down the stress level.
func (logger *Logger) Info(methodType string, trace ...string) Writer {
	if logger.Level < LevelInfo {
		return Empty{}
	}
	pc, file, line, ok := runtime.Caller(2)
	method := "unknown"
	if ok {
		details := runtime.FuncForPC(pc)
		method = filepath.Base(details.Name())
	}
	return &Content{
		isTerm: logger.isTerm,
		logger: logger,
		Type:   methodType,
		Method: method,
		Level:  LevelInfo,
		File:   file,
		Line:   line,
		Traces: trace,
	}
}

// Debug With DEBUG, you start to include more granular, diagnostic information.
func Debug(methodType string, trace ...string) Writer {
	return DefaultLogger.Debug(methodType, trace...)
}

// Debug With DEBUG, you start to include more granular, diagnostic information.
func (logger *Logger) Debug(methodType string, trace ...string) Writer {
	if logger.Level < LevelDebug {
		return Empty{}
	}
	pc, file, line, ok := runtime.Caller(2)
	method := "unknown"
	if ok {
		details := runtime.FuncForPC(pc)
		method = filepath.Base(details.Name())
	}
	return &Content{
		isTerm: logger.isTerm,
		logger: logger,
		Type:   methodType,
		Method: method,
		Level:  LevelDebug,
		File:   file,
		Line:   line,
		Traces: trace,
	}
}
