package logger

import (
	"fmt"
	"net/http"
	"time"
)

// Trace set the traces for content
func (content *Content) Trace(trace ...string) *Content {
	content.Traces = append(content.Traces, trace...)
	return content
}

// Message set the message for content
func (content *Content) Message(message string) *Content {
	content.Msg = message
	return content
}

// Write write log to the output writer
func (content *Content) Write() {
	if content.logger.Level < content.Level {
		return
	}
	content.TimeStamp = time.Now()
	if content.StartAt != nil {
		content.Latency = content.TimeStamp.Sub(*content.StartAt)
	}
	log := content.logger.Formatter(*content)
	fmt.Fprintln(content.logger.Output, log)
}

// TimeCost count the time costs and  write log to the output writer
func (content *Content) TimeCost(start time.Time) {
	content.StartAt = &start
	content.Write()
}

// CodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (content *Content) CodeColor() string {
	code := content.Code
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
func (content *Content) MethodColor() string {
	switch content.Type {
	case CREATE:
		return blue
	case UPDATE:
		return yellow
	case RETRIEVE:
		return green
	case DELETE:
		return magenta
	case ERROR:
		return red
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (content *Content) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (content *Content) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && content.isTerm)
}
