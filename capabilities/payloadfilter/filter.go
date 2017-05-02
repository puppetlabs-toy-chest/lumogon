package payloadfilter

import (
	"fmt"
)

// Filter takes a map intended to be passed along as a capability.Payload, and
// pre-processes the map data to ensure that it is valid for use. This includes
// removing any empty strings, lists, and maps, while ensuring that something is
// available in the top-level map.
func Filter(input map[string]interface{}) (map[string]interface{}, error) {
	result, err := filterMap(input)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return map[string]interface{}{}, nil
	}
	return result, err
}

func filterMap(value map[string]interface{}) (map[string]interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}

	result := make(map[string]interface{})

	for key, val := range value {
		if key == "" {
			continue
		}

		switch coerced := val.(type) {
		case string:
			if coerced != "" {
				result[key] = coerced
			}

		case map[string]interface{}:
			nested, err := filterMap(coerced)
			if err != nil {
				return nil, err
			}
			if nested != nil {
				result[key] = nested
			}

		case map[string]string:
			nested := make(map[string]string)
			for k, v := range coerced {
				if k == "" || v == "" {
					continue
				}
				nested[k] = v
			}
			if len(nested) != 0 {
				result[key] = nested
			}
		default:
			err := fmt.Errorf("Payload contains data which is not a string, map of strings, or map of maps: %#+v", coerced)
			return nil, err
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
