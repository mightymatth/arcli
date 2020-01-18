package utils

import (
	"fmt"
	"strings"
)

func PrintWithDelimiter(s []string) string {
	newS := make([]string, 0, len(s))

	for _, val := range s {
		newS = append(newS, fmt.Sprintf("'%s'", val))
	}

	return strings.Join(newS, ", ")
}
