package logger

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	defer LogR("Message", time.Now())
}
