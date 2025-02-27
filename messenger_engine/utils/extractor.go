package utils

import "fmt"

// ExtractInt safely extracts an integer value from a map.
func ExtractInt(data map[string]interface{}, key string) (int, error) {
	value, ok := data[key].(float64) // JSON numbers are float64
	if !ok {
		return 0, fmt.Errorf("%s is missing or invalid", key)
	}
	return int(value), nil
}

// ExtractString safely extracts a string value from a map.
func ExtractString(data map[string]interface{}, key string) (string, error) {
	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("%s is missing or invalid", key)
	}
	return value, nil
}

// ExtractBool safely extracts a boolean value from a map, with a default fallback.
func ExtractBool(data map[string]interface{}, key string, defaultValue bool) bool {
	value, ok := data[key].(bool)
	if !ok {
		return defaultValue
	}
	return value
}
