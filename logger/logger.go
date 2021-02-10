// Fork from https://github.com/gin-gonic/gin/blob/master/logger.go

package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/yaoapp/xun/global"
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

// Method Type definition, including "CREATE","UPDATE","RETRIEVE" and "DELETE" .
const (
	MethodC = "CREATE"
	MethodU = "UPDATE"
	MethodR = "RETRIEVE"
	MethodD = "DELETE"
)

var consoleColorMode = autoColor

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	// TimeStamp shows the time after the method is executed.
	TimeStamp time.Time
	// Latency is how much time the method cost to process a certain request.
	Latency time.Duration
	// MethodType is the type of method, including "CREATE","UPDATE","RETRIEVE" and "DELETE" .
	MethodType string
	// Method is the executed method name.
	Method string
	// SQL is executed SQL.
	SQL string
	// StatusCode is a code method response, same as  HTTP response code.
	StatusCode int
	// Message is set if error has occurred in processing the request.
	Message string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// LogFormatter gives the signature of the formatter function passed to WithFormatter
type LogFormatter func(params LogFormatterParams) string

// Config defines the config for Logger middleware.
type Config struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipMethods is a method array which logs are not written.
	// Optional.
	SkipMethods []string
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	switch p.MethodType {
	case MethodC:
		return yellow
	case MethodU:
		return magenta
	case MethodR:
		return green
	case MethodD:
		return magenta
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[XUN] %v |%s %3d %s| %13v |%s %-8s %s %#v\n%s\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		methodColor, param.MethodType, resetColor,
		param.Method,
		param.Message,
		param.SQL,
	)
}

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	consoleColorMode = disableColor
}

// ForceConsoleColor force color output in the console.
func ForceConsoleColor() {
	consoleColorMode = forceColor
}

// HandlerFunc the handler of logger
type HandlerFunc func(string)

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() HandlerFunc {
	return WithConfig(Config{})
}

// WithFormatter instance a Logger middleware with the specified log format function.
func WithFormatter(f LogFormatter) HandlerFunc {
	return WithConfig(Config{
		Formatter: f,
	})
}

// WithWriter instance a Logger middleware with the specified writer buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func WithWriter(out io.Writer, notlogged ...string) HandlerFunc {
	return WithConfig(Config{
		Output:      out,
		SkipMethods: notlogged,
	})
}

// WithConfig instance a Logger middleware with config.
func WithConfig(conf Config) HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = global.DefaultWriter
	}

	notlogged := conf.SkipMethods
	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(called string) {
		// Start timer
		start := time.Now()
		method := called
		// log only when path is not being skipped
		if _, ok := skip[method]; !ok {
			param := LogFormatterParams{
				isTerm:     isTerm,
				Keys:       map[string]interface{}{"xx": "xxxx"},
				SQL:        "SQL",
				MethodType: MethodC,
				Method:     method,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)
			param.StatusCode = 200
			param.Message = "Message"

			fmt.Fprint(out, formatter(param))
		}
	}
}

// Debug debug out
func Debug(sql string, i ...interface{}) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok {
		// method := fmt.Sprintf("called from %s#%d\n", details, no)
		Logger()(details.Name())
	}
}

// LogC logging
func LogC(i ...interface{}) {}

// LogU logging
func LogU(i ...interface{}) {}

// LogR logging
func LogR(i ...interface{}) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok {
		// method := fmt.Sprintf("called from %s#%d\n", details, no)
		Logger()(details.Name())
	}
}

// LogD logging
func LogD(i ...interface{}) {}
