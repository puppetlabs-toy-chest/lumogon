package utils

import "time"

// GetTimestamp returns the current runtime
// TODO - need to decide the best format for this
func GetTimestamp() string {
	return time.Now().String()
}
