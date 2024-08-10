package utils

import "strings"

func KeyValueStringsToMap(values []string) map[string]string {
	result := make(map[string]string, len(values))

	for _, value := range values {
		kv := strings.SplitN(value, "=", 2) //nolint:mnd
		if len(kv) == 1 {
			result[kv[0]] = ""
		} else {
			result[kv[0]] = kv[1]
		}
	}

	return result
}
