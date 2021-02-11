package logger

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
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

// Method Type definition, including "CREATE","UPDATE","RETRIEVE","DELETE" and "ERROR".
const (
	CREATE   = "CREATE"
	UPDATE   = "UPDATE"
	RETRIEVE = "RETRIEVE"
	DELETE   = "DELETE"
	ERROR    = "ERROR"
)

// LogLevel the log levels definition.
type LogLevel int

// the log levels
const (
	LevelFatal LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
	LevelAll
)

var levelNames = map[LogLevel]string{
	LevelFatal: "FATAL",
	LevelError: "ERROR",
	LevelWarn:  "WARN",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
	LevelTrace: "TRACE",
}

var consoleColorMode = autoColor

// Content is the structure any formatter will be handed when time to log comes
type Content struct {
	// TimeStamp shows the time after the method is executed.
	TimeStamp time.Time
	// StartAt shows the time the caller method is executed.
	StartAt *time.Time
	// Latency is how much time the method cost to process a certain request.
	Latency time.Duration
	// The level of log
	Level LogLevel
	// Type is the type of the log, including "CREATE","UPDATE","RETRIEVE" and "DELETE" .
	Type string
	// Method is the executed method name.
	Method string
	// Traces is executed SQL list.
	Traces []string
	// Code is a code method response, same as  HTTP response code.
	Code int
	// Message is set if error has occurred in processing the request.
	Msg string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// File the file where the method executed is
	File string
	// Line the line number where the method executed is
	Line int
	// Contexts are the Context set on the request's context.
	Contexts interface{}

	// The logger intance
	logger *Logger
}

// Formatter gives the signature of the formatter function passed to WithFormatter
type Formatter func(content Content) string

// Logger the logger instance
type Logger struct {

	// Optional. Default value is defaultLogFormatter
	Formatter Formatter

	// Output is a writer where logs are written.
	// Optional. Default value is xun.DefaultWriter.
	Output io.Writer

	// SkipMethods is a method array which logs are not written.
	// Optional.
	SkipMethods []string

	// IsTerm shows whether does xun's output descriptor refers to a terminal.
	isTerm bool

	// Level display witch level, default is all
	Level LogLevel
}

// DefaultLogger the default logger instance
var DefaultLogger *Logger

// DefaultErrorLogger the default error logger instance
var DefaultErrorLogger *Logger

// DefaultFormatter is the default log format function Logger middleware uses.
var DefaultFormatter = func(c Content) string {
	var methodColor, resetColor string
	if c.IsOutputColor() {
		// statusColor = param.CodeColor()
		methodColor = c.MethodColor()
		resetColor = c.ResetColor()
	}
	if c.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		c.Latency = c.Latency - c.Latency%time.Second
	}

	traces := strings.Join(c.Traces, "\n")
	if traces != "" {
		traces = "\n" + traces
	}

	message := c.Msg
	if message != "" {
		message = fmt.Sprintf("\n%s %s %s", methodColor, c.Msg, resetColor)
	}

	level := levelNames[c.Level]

	if c.Type == ERROR {
		return fmt.Sprintf("[XUN] %v | %-5s |%s %3d %s| %13v | %s %s %s",
			c.TimeStamp.Format("2006/01/02 - 15:04:05"),
			level,
			methodColor, c.Code, resetColor,
			c.Latency,
			c.Method,
			message,
			traces,
		)
	}
	return fmt.Sprintf("[XUN] %v | %-5s |%s %-8s %s| %13v | %s %s %s",
		c.TimeStamp.Format("2006/01/02 - 15:04:05"),
		level,
		methodColor, c.Type, resetColor,
		c.Latency,
		c.Method,
		message,
		traces,
	)
}

// DefaultLevel the default level
var DefaultLevel = LevelAll
