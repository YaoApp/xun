//  Fork from https://github.com/gin-gonic/gin/blob/master/mode.go

package global

import (
	"io"
	"os"
)

// EnvXunMode indicates environment name for xun mode.
const EnvXunMode = "XUN_MODE"

const (
	// DebugMode indicates xun mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates xun mode is release.
	ReleaseMode = "release"
	// TestMode indicates xun mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

// DefaultWriter is the default io.Writer used by Xun for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
// 		import "github.com/mattn/go-colorable"
// 		xun.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Xun to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

var xunMode = debugCode
var modeName = DebugMode

func init() {
	mode := os.Getenv(EnvXunMode)
	SetMode(mode)
}

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	if value == "" {
		value = DebugMode
	}

	switch value {
	case DebugMode:
		xunMode = debugCode
	case ReleaseMode:
		xunMode = releaseCode
	case TestMode:
		xunMode = testCode
	default:
		panic("xun mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}

// Mode returns currently xun mode.
func Mode() string {
	return modeName
}
