package utils

import (
	"fmt"
	"strings"
)

// PrintWithDelimiter returns string that contains array of single-quoted
// strings delimited with a comma.
func PrintWithDelimiter(s []string) string {
	newS := make([]string, 0, len(s))

	for _, val := range s {
		newS = append(newS, fmt.Sprintf("'%s'", val))
	}

	return strings.Join(newS, ", ")
}
