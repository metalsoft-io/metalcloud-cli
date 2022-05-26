package main

import (
	"fmt"
	"strings"
)

func truncateString(s string, length int) string {

	if len(s) <= length {
		return s
	}

	str := s
	if len(str) > 0 {
		return str[:length] + "..."
	}
	return ""
}

//makes multi-line string from very long string
func wrapToLength(s string, length int) string {

	if len(s) <= length {
		return s
	}

	rows := int(len(s) / length)

	var sb strings.Builder

	for i := 0; i < rows; i++ {
		sb.WriteString(fmt.Sprintf("%s\n", s[i*length:(i+1)*length]))
	}
	remain := int(len(s) % length)
	sb.WriteString(fmt.Sprintf("%s", s[len(s)-remain:]))

	return sb.String()
}
