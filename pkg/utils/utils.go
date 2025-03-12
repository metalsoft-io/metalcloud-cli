package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func CreateSlug(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}

	processedString := reg.ReplaceAllString(input, " ")

	processedString = strings.TrimSpace(processedString)

	slug := strings.ReplaceAll(processedString, " ", "-")

	slug = strings.ToLower(slug)

	return slug
}

func ReadConfigFromPipe() ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no data piped to stdin")
	}

	content, err := os.ReadFile("/dev/stdin")
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ReadConfigFromFile(configSource string) ([]byte, error) {
	content, err := os.ReadFile(configSource)
	if err != nil {
		return nil, err
	}

	return content, nil
}
