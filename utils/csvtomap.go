package utils

import (
	"fmt"
	"strings"
)

// CsvToMap TODO
func CsvToMap(csv []string) (map[string]string, error) {
	var err error
	output := make(map[string]string)
	if len(csv) == 0 {
		return output, nil
	}
	for _, entry := range csv {
		keyvalue := strings.Split(entry, ",")
		if len(keyvalue) != 2 {
			err = fmt.Errorf("Error splitting entry, expected a key and value but extracted %d elements: %v", len(keyvalue), keyvalue)
			break
		}
		if keyvalue[0] == "" {
			err = fmt.Errorf("Error splitting entry, empty key detected in entry: %s", entry)
			break
		}
		if keyvalue[1] == "" {
			err = fmt.Errorf("Error splitting entry, empty value detected in entry: %s", entry)
			break
		}
		output[keyvalue[0]] = keyvalue[1]
	}
	return output, err
}
