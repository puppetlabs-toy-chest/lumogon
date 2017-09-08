package versions

import (
	"github.com/docker/docker/api/types/versions"
)

// LessThan checks if a version is less than another
func LessThan(v, other string) bool {
	return versions.LessThan(v, other)
}
