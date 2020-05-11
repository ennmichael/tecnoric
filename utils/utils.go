package utils

import "strings"

func SplitAndTrim(s, sep string) []string {
	var result []string
	for _, ss := range strings.Split(s, sep) {
		result = append(result, strings.TrimSpace(ss))
	}
	return result
}

