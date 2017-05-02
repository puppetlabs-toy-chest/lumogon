package logging

import (
	"log"
	"os"
)

var (
	// Debug flag
	Debug bool
)

var (
	stderr *log.Logger
	stdout *log.Logger
)

// Stdout provides a convenient way to log to STDOUT
func Stdout(format string, args ...interface{}) {
	if args != nil {
		stderr.Printf(format, args...)
	} else {
		stderr.Println(format)
	}
}

// Stderr provides a convenient way to log to STDERR
func Stderr(format string, args ...interface{}) {
	if !Debug {
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
	stderr = log.New(os.Stderr, "[lumogon] ", logProperties)
	stdout = log.New(os.Stdout, "", logProperties)
}
