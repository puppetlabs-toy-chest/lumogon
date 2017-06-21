package logging

import (
	"log"
	"os"
)

var (
	// DebugEnabled flag
	DebugEnabled bool
)

var (
	stderr *log.Logger
)

// Info provides a convenient way to log to STDERR
// Note that we use STDERR rather than STDOUT as the primary
// purpose of this logging interface is "diagnostic messages"
// as per POSIX
func Info(format string, args ...interface{}) {
	if args != nil {
		stderr.Printf(format, args...)
	} else {
		stderr.Println(format)
	}
}

// Debug provides a convenient way to log to STDERR only when the
// debug flag is passed
func Debug(format string, args ...interface{}) {
	if !DebugEnabled {
		return
	}
	if args != nil {
		stderr.Printf(format, args...)
	} else {
		stderr.Println(format)
	}
}

func init() {
	logProperties := log.Ldate | log.Ltime | log.Lmicroseconds
	stderr = log.New(os.Stderr, "[Lumogon] ", logProperties)
}
