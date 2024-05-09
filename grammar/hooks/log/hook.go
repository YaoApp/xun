package log

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/qustavo/sqlhooks/v2"
	"github.com/yaoapp/kun/log"
	"github.com/yaoapp/xun/grammar/hooks"
)

func init() {
	err := hooks.RegisterHook("log", Default)
	if err != nil {
		panic(err)
	}
}

// Default default log hook instance
var Default = &Hook{
	Logger:       &KunLogger{},
	Level:        log.InfoLevel,
	MaxFieldSize: 256,
}

var (
	ErrorFieldName       = "error"
	QueryFieldName       = "query"
	RequestTimeFieldName = "rt"
	ArgFieldPrefix       = "arg_"
)

// Hook record sql logs
type Hook struct {
	Level         log.Level `json:"level,omitempty"`
	MaxFieldSize  int       `json:"max_field_size,omitempty"`
	Logger        Logger
	ContextFields func(ctx context.Context) log.F
}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, ctxKeyStartTime, time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return h.log(ctx, query, nil, args...)
}

func (h *Hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	if errors.Is(err, driver.ErrSkip) {
		return err
	}
	_, _ = h.log(ctx, query, err, args...)
	return err
}

func (h *Hook) log(ctx context.Context, query string, err error, args ...interface{}) (context.Context, error) {
	start := ctx.Value(ctxKeyStartTime).(time.Time)
	rt := time.Since(start).Nanoseconds() / 1e6 // unit: Ms
	fields := make(log.F)
	if err != nil {
		fields[ErrorFieldName] = err
	}
	fields[QueryFieldName] = query
	fields[RequestTimeFieldName] = rt
	for k, v := range h.ContextFields(ctx) {
		fields[k] = v
	}
	for i, arg := range args {
		argName := ArgFieldPrefix + strconv.Itoa(i)
		if h.MaxFieldSize-ellipsisLength <= 0 {
			fields[argName] = arg
		} else {
			fields[argName] = LimitSize(arg, h.MaxFieldSize-ellipsisLength)
		}
	}
	h.Logger.Log(h.Level, "sql log", fields)
	return ctx, nil
}

// LimitSize limit the print size of the value
func LimitSize(value interface{}, n int) interface{} {
	switch val := value.(type) {
	case []bool, []complex128, []complex64, []float64, []float32,
		[]int, []int64, []int32, []int16, []int8, []string, []uint, []uint64, []uint32, []uint16, []uintptr, []time.Time,
		[]time.Duration, []error:
		s := fmt.Sprintf("%v", value)
		if len(s) > n {
			return s[:n] + "..."
		}
		return value
	case string:
		if len(val) > n {
			return val[:n] + "..."
		}
		return val
	case *string:
		if val == nil {
			return "<nil>"
		}
		if len(*val) > n {
			return (*val)[:n] + "..."
		}
		return *val
	case []byte:
		if val == nil {
			return []byte("<nil>")
		}
		if len(val) > n {
			return string(append(val[:n], []byte("...")...))
		}
		return string(val)
	default:
		return value
	}
}

var (
	ctxKeyStartTime = struct{}{}
	ellipsisLength  = 3
)

var _ sqlhooks.Hooks = (*Hook)(nil)
