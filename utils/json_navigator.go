package utils

import (
	"fmt"
	"handbook-scraper/utils/log"
	"reflect"
	"strconv"
)

// GetTypedValue retrieves a value from a JSON map and attempts to cast it to the specified type.
// If the value is not found or cannot be cast, it returns the default value of type T.
func GetTypedValue[T any](data map[string]interface{}, path string) T {
	var zero T // Default zero value for type T

	// Find the interface value at the specified path
	value, err := findInterface(data, path)
	if err != nil {
		log.Errorf("Error retrieving value at path '%s': %v\n", path, err)
		return zero
	}

	// Handle nil values
	if value == nil {
		log.Errorf("Value at path '%s' is nil\n", path)
		return zero
	}

	// Handle []map[string]interface{} specifically
	if _, isTargetType := any(zero).([]map[string]interface{}); isTargetType {
		if slice, ok := value.([]interface{}); ok {
			var result []map[string]interface{}
			for _, item := range slice {
				if itemMap, ok := item.(map[string]interface{}); ok {
					result = append(result, itemMap)
				} else {
					log.Warnf("Item in slice at path '%s' is not a map[string]interface{}\n", path)
				}
			}
			return any(result).(T)
		}
		log.Warnf("Value at path '%s' is not a slice of interface{} for []map[string]interface{}\n", path)
		return zero
	}

	// Handle string to numeric type conversion
	switch any(zero).(type) {
	case int:
		if str, ok := value.(string); ok {
			if intValue, err := strconv.Atoi(str); err == nil {
				return any(intValue).(T)
			} else {
				log.Warnf("Failed to convert string '%s' to int at path '%s': %v\n", str, path, err)
			}
		} else {
			log.Warnf("Value at path '%s' is not a string for int conversion\n", path)
		}
	case float32:
		if str, ok := value.(string); ok {
			if floatValue, err := strconv.ParseFloat(str, 32); err == nil {
				return any(float32(floatValue)).(T)
			} else {
				log.Warnf("Failed to convert string '%s' to float32 at path '%s': %v\n", str, path, err)
			}
		} else {
			log.Warnf("Value at path '%s' is not a string for float32 conversion\n", path)
		}
	case float64:
		if str, ok := value.(string); ok {
			if floatValue, err := strconv.ParseFloat(str, 64); err == nil {
				return any(floatValue).(T)
			} else {
				log.Warnf("Failed to convert string '%s' to float64 at path '%s': %v\n", str, path, err)
			}
		} else {
			log.Warnf("Value at path '%s' is not a string for float64 conversion\n", path)
		}
	case bool:
		if str, ok := value.(string); ok {
			if boolValue, err := strconv.ParseBool(str); err == nil {
				return any(boolValue).(T)
			} else {
				log.Warnf("Failed to convert string '%s' to bool at path '%s': %v\n", str, path, err)
			}
		} else {
			log.Warnf("Value at path '%s' is not a string for bool conversion\n", path)
		}
	}

	// Attempt to assert the value to the desired type T
	if typedValue, ok := value.(T); ok {
		return typedValue
	}

	// Handle custom types with type conversion using reflection
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		log.Warnf("Value at path '%s' is invalid for reflection\n", path)
		return zero
	}

	zeroType := reflect.TypeOf(zero)
	if val.Type().ConvertibleTo(zeroType) {
		convertedValue := val.Convert(zeroType).Interface()
		if convertedTypedValue, ok := convertedValue.(T); ok {
			return convertedTypedValue
		}
	}

	log.Warnf("Value at path '%s' is not of type %T\n", path, zero)
	return zero
}

// findInterface navigates the JSON map using the provided path and returns the value as an interface{}.
func findInterface(data map[string]interface{}, path string) (interface{}, error) {
	current := data
	keys := splitFieldPath(path)
	var traversedPath []string

	for _, key := range keys {
		traversedPath = append(traversedPath, key)
		value, exists := current[key]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found at '%s'", key, joinPath(traversedPath))
		}

		// Return value if we're at the last key
		if len(traversedPath) == len(keys) {
			return value, nil
		}

		// Check if value is a nested map
		nestedMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field '%s' is not a nested object at '%s'", key, joinPath(traversedPath))
		}
		current = nestedMap
	}

	return nil, fmt.Errorf("invalid path '%s'", path)
}

// splitFieldPath splits a dot-separated field path into its components.
// For example, "props.pageProps.handbook_synopsis" becomes ["props", "pageProps", "handbook_synopsis"].
func splitFieldPath(path string) []string {
	var keys []string
	current := ""
	for _, char := range path {
		if char == '.' {
			if current != "" {
				keys = append(keys, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		keys = append(keys, current)
	}
	return keys
}

// joinPath joins the keys back into a dot-separated path for error messages.
// For example, ["props", "pageProps", "handbook_synopsis"] becomes "props.pageProps.handbook_synopsis".
func joinPath(keys []string) string {
	path := ""
	for i, key := range keys {
		if i > 0 {
			path += "."
		}
		path += key
	}
	return path
}
