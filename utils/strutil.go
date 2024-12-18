package utils

import "strings"

func KeyValueStringsToMap(values []string) map[string]string {
	result := make(map[string]string, len(values))

	for _, value := range values {
		k, v, _ := strings.Cut(value, "=")
		result[k] = v
	}

	return result
}
