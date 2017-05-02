package utils

import "runtime"

// CurrentFunctionName returns the scope of the function from which it is invoked
func CurrentFunctionName() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
