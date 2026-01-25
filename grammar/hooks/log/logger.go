package log

import (
	"github.com/yaoapp/kun/log"
)

// Logger log hook logger interface
type Logger interface {
	Log(level log.Level, msg string, fields log.F)
}

// KunLogger implements Logger with kun/log
type KunLogger struct{}

func (l *KunLogger) Log(level log.Level, msg string, fields log.F) {
	switch level {
	case log.TraceLevel:
		log.With(fields).Trace("%s", msg)
	case log.DebugLevel:
		log.With(fields).Debug("%s", msg)
	case log.InfoLevel:
		log.With(fields).Info("%s", msg)
	case log.WarnLevel:
		log.With(fields).Warn("%s", msg)
	case log.ErrorLevel:
		log.With(fields).Error("%s", msg)
	case log.PanicLevel:
		log.With(fields).Panic("%s", msg)
	case log.FatalLevel:
		log.With(fields).Fatal("%s", msg)
	}
}

var _ Logger = (*KunLogger)(nil)
