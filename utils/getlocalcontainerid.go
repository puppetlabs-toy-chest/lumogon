package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// GetLocalContainerID returns the ID of the current local container
func GetLocalContainerID(cgroupfile string) (string, error) {
	file, err := os.Open(cgroupfile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	re := regexp.MustCompile(`^\d+:cpu:/docker/(\S+)$`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) != 0 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("No containerID found in %s", cgroupfile)
}
