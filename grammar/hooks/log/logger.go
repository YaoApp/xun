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
		log.With(fields).Trace(msg)
	case log.DebugLevel:
		log.With(fields).Debug(msg)
	case log.InfoLevel:
		log.With(fields).Info(msg)
	case log.WarnLevel:
		log.With(fields).Warn(msg)
	case log.ErrorLevel:
		log.With(fields).Error(msg)
	case log.PanicLevel:
		log.With(fields).Panic(msg)
	case log.FatalLevel:
		log.With(fields).Fatal(msg)
	}
}

var _ Logger = (*KunLogger)(nil)
