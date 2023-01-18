package main

import (
	"fmt"
	"strings"
	"unicode"
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

// makes multi-line string from very long string
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

// returns a label compatible string from any string
// will throw errors if the string starts with numbers or -
// will truncate all other non alpha and dash chars
// will throw error if remaining string is empty
func makeLabel(s string) (string, error) {

	if s == "" {
		return "", fmt.Errorf("label cannot be empty")
	}

	label := ""

	for i, r := range strings.ToLower(s) {
		if i == 0 {
			if !unicode.IsLetter(r) {
				return "", fmt.Errorf("label must start with a letter ")
			}
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == rune('-') {
			label = label + string(r)
		}
	}

	if label == "" {
		return "", fmt.Errorf("generated label is empty. This means that the letters in the provided name are not alphanumeric")
	}

	return label, nil
}
