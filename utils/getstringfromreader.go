package utils

import (
	"bytes"
	"fmt"
	"io"
)

// GetStringFromReader converts io.Reader to string
func GetStringFromReader(reader io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return "", fmt.Errorf("Unable to convert Reader to string, error: %s", err)
	}

	return buf.String(), nil
}
