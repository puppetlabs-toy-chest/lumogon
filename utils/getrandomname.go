package utils

import (
	"fmt"

	"github.com/docker/docker/pkg/namesgenerator"
)

// GetRandomName returns a random name prefixed with the supplied string
func GetRandomName(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, namesgenerator.GetRandomName(1))
}
